param(
    [string]$TestDir = ""
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($TestDir)) {
    $TestDir = $PSScriptRoot
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")

$schemaIntent  = Join-Path $repoRoot "schemas\intent-envelope.schema.json"
$schemaAdm     = Join-Path $repoRoot "schemas\admission_result_v0.1.schema.json"
$schemaCmd     = Join-Path $repoRoot "schemas\command-envelope.schema.json"
$regCapability = Join-Path $repoRoot "core-capability-registry.yaml"
$regAggregate  = Join-Path $repoRoot "aggregate-ownership-registry.yaml"
$regCommand    = Join-Path $repoRoot "command-event-catalogue.yaml"
$regBaseline   = Join-Path $repoRoot "architecture-baseline.yaml"

function Read-TextFile($path) {
    return (Get-Content -Path $path -Raw) -replace "`r`n", "`n"
}

function Parse-Capabilities($path) {
    $content = Read-TextFile $path
    $caps = @{}
    [regex]::Matches($content, '(?ms)^  - id: ([^\n]+)\n[\s\S]*?^    mutation: (true|false)') | ForEach-Object {
        $id = $_.Groups[1].Value.Trim()
        $mutation = $_.Groups[2].Value.Trim() -eq 'true'
        $caps[$id] = $mutation
    }
    return $caps
}

function Parse-Ownership($path) {
    $content = Read-TextFile $path
    $owned = @{}
    [regex]::Matches($content, '(?ms)^  - aggregate_id: ([^\n]+)\n[\s\S]*?^    owned_capabilities:\n([\s\S]*?)^    description:') | ForEach-Object {
        $agg = $_.Groups[1].Value.Trim()
        $list = $_.Groups[2].Value -split '\n' | ForEach-Object { if ($_ -match '^      - (.+)') { $matches[1].Trim() } } | Where-Object { $_ }
        $owned[$agg] = $list
    }
    return $owned
}

function Parse-Commands($path) {
    $content = Read-TextFile $path
    $commands = @{}
    [regex]::Matches($content, '(?ms)^  - command_id: ([^\n]+)\n[\s\S]*?^    transition_id: ([^\n]+)') | ForEach-Object {
        $block = $_.Value
        $cmd = [regex]::Match($block, '^  - command_id: ([^\n]+)').Groups[1].Value.Trim()
        $cap = [regex]::Match($block, '\n    capability_id: ([^\n]+)').Groups[1].Value.Trim()
        $agg = [regex]::Match($block, '\n    aggregate_id: ([^\n]+)').Groups[1].Value.Trim()
        $trans = [regex]::Match($block, '\n    transition_id: ([^\n]+)').Groups[1].Value.Trim()
        $commands[$cmd] = @{
            capability_id = $cap
            aggregate_id  = $agg
            transition_id = $trans
        }
    }
    return $commands
}

function Get-BaselineRevision($path) {
    $content = Read-TextFile $path
    $m = [regex]::Match($content, '(?m)^  revision: "?([^"\n]+)"?')
    if (-not $m.Success) { throw "Could not find architecture baseline revision" }
    return $m.Groups[1].Value.Trim()
}

function Validate-Json($path, $schema) {
    $valid = $false
    $err = $null
    try { $valid = Test-Json -Path $path -SchemaFile $schema } catch { $err = $_.Exception.Message }
    return $valid, $err
}

function New-AdmissionResult($admissionId, $decision, $capRef, $aggRef, $cmdRef, $transRef, $reasons, $rules) {
    return [ordered]@{
        schema_version        = "v0.1"
        architecture_revision = $baselineRevision
        admission_id          = $admissionId
        decision              = $decision
        capability_ref        = $capRef
        aggregate_ref         = $aggRef
        command_ref           = $cmdRef
        transition_ref        = $transRef
        rule_refs             = $rules
        reason_codes          = $reasons
    }
}

function Get-PropertyOrDefault($obj, $name, $default) {
    if ($obj -and $obj.PSObject.Properties[$name]) { return $obj.$name }
    return $default
}

$capabilities     = Parse-Capabilities $regCapability
$ownership        = Parse-Ownership $regAggregate
$commands         = Parse-Commands $regCommand
$baselineRevision = Get-BaselineRevision $regBaseline

$caseDirs = Get-ChildItem -Directory -Path $TestDir
$failures = 0

foreach ($case in $caseDirs) {
    $inputFile    = Join-Path $case.FullName "input.json"
    $expectedFile = Join-Path $case.FullName "expected.json"

    if (-not (Test-Path $inputFile) -or -not (Test-Path $expectedFile)) {
        continue
    }

    $intent = $null
    try { $intent = Get-Content $inputFile | ConvertFrom-Json } catch { }
    if (-not $intent) {
        Write-Output ("{0}: INVALID JSON input" -f $case.Name)
        $failures++
        continue
    }

    $validIntent, $intentErr = Validate-Json $inputFile $schemaIntent
    $reasons = @()
    $rules   = @()
    $decision = "ADMIT"
    $capRef   = Get-PropertyOrDefault $intent 'capability_ref' ''
    $aggRef   = Get-PropertyOrDefault $intent 'aggregate_ref'  ''
    $cmdRef   = Get-PropertyOrDefault $intent 'command_ref'    ''
    $transRef = ''
    $admissionId = "rs-001-{0}" -f $case.Name

    if (-not $validIntent) {
        if ($intentErr -match 'Required properties \["([^"]+)"\]') {
            $reasons += "MISSING_REQUIRED_REFERENCE"
        } else {
            $reasons += "INVALID_INTENT"
        }
        $decision = "REJECT"
    } else {
        # IR-01: architecture baseline
        if ($intent.architecture_revision -ne $baselineRevision) {
            $reasons += "ARCHITECTURE_REVISION_MISMATCH"
            $rules   += "IR-01"
        } else {
            $rules   += "IR-01"
        }

        # IR-03: capability exists and is mutating
        if (-not $reasons) {
            if (-not $capabilities.ContainsKey($capRef)) {
                $reasons += "UNKNOWN_CAPABILITY"
                $rules   += "IR-03"
            } elseif (-not $capabilities[$capRef]) {
                $reasons += "CAPABILITY_NOT_MUTATING"
                $rules   += "IR-03"
            } else {
                $rules   += "IR-03"
            }
        }

        # IR-04: aggregate ownership
        if (-not $reasons) {
            $owner = $null
            foreach ($agg in $ownership.Keys) {
                if ($ownership[$agg] -contains $capRef) {
                    $owner = $agg
                    break
                }
            }
            if ($owner -ne $aggRef) {
                $reasons += "OWNERSHIP_MISMATCH"
                $rules   += "IR-04"
            } else {
                $rules   += "IR-04"
            }
        }

        # IR-05: command catalogue and mapping
        if (-not $reasons) {
            if (-not $commands.ContainsKey($cmdRef)) {
                $reasons += "UNKNOWN_COMMAND"
                $rules   += "IR-05"
            } else {
                $cmd = $commands[$cmdRef]
                if ($cmd.capability_id -ne $capRef) {
                    $reasons += "COMMAND_CAPABILITY_MISMATCH"
                    $rules   += "IR-05"
                } elseif ($cmd.aggregate_id -ne $aggRef) {
                    $reasons += "COMMAND_AGGREGATE_MISMATCH"
                    $rules   += "IR-05"
                } elseif ([string]::IsNullOrWhiteSpace($cmd.transition_id)) {
                    $reasons += "UNDEFINED_TRANSITION"
                    $rules   += "IR-05"
                } else {
                    $transRef = $cmd.transition_id
                    $rules   += "IR-05"
                }
            }
        }

        if ($reasons) {
            $decision = "REJECT"
        } else {
            $rules += "IR-10"
        }
    }

    if ([string]::IsNullOrWhiteSpace($capRef))   { $capRef   = "not-provided" }
    if ([string]::IsNullOrWhiteSpace($aggRef))   { $aggRef   = "not-provided" }
    if ([string]::IsNullOrWhiteSpace($cmdRef))   { $cmdRef   = "not-provided" }
    if ([string]::IsNullOrWhiteSpace($transRef)) { $transRef = "not-provided" }
    if ($rules.Count -eq 0) { $rules = @("IR-03", "IR-04", "IR-05", "IR-10") }
    if ($decision -eq "ADMIT" -and $rules.Count -eq 0) { $rules = @("IR-03", "IR-04", "IR-05", "IR-10") }

    $result = New-AdmissionResult $admissionId $decision $capRef $aggRef $cmdRef $transRef $reasons $rules
    $resultPath = Join-Path $case.FullName "actual_admission_result.json"
    $result | ConvertTo-Json -Depth 10 | Set-Content $resultPath

    $validAdm, $errAdm = Validate-Json $resultPath $schemaAdm
    if (-not $validAdm) {
        Write-Output ("{0}: FAIL actual_admission_result schema: {1}" -f $case.Name, $errAdm)
        $failures++
        continue
    }

    $expected = Get-Content $expectedFile | ConvertFrom-Json
    $reasonMatch = ($result.reason_codes | Sort-Object) -join ',' -eq ($expected.reason_codes | Sort-Object) -join ','

    if ($result.decision -ne $expected.decision -or -not $reasonMatch) {
        Write-Output ("{0}: FAIL expected decision={1} reason_codes={2}, got decision={3} reason_codes={4}" -f $case.Name, $expected.decision, ($expected.reason_codes -join ','), $result.decision, ($result.reason_codes -join ','))
        $failures++
        continue
    }

    if ($decision -eq "ADMIT") {
        $cmdEnv = [ordered]@{
            schema_version        = "v0.1"
            architecture_revision = $baselineRevision
            command_id            = "rs-001-cmd-{0}" -f $case.Name
            admission_ref         = $admissionId
            capability_ref        = $capRef
            aggregate_ref         = $aggRef
            command_ref           = $cmdRef
            transition_ref        = $transRef
            payload               = $intent.payload
        }
        $cmdEnvPath = Join-Path $case.FullName "actual_command_envelope.json"
        $cmdEnv | ConvertTo-Json -Depth 10 | Set-Content $cmdEnvPath

        $validCmd, $errCmd = Validate-Json $cmdEnvPath $schemaCmd
        if (-not $validCmd) {
            Write-Output ("{0}: FAIL command_envelope schema: {1}" -f $case.Name, $errCmd)
            $failures++
            continue
        }
    }

    Write-Output ("{0}: PASS (decision={1}, reason_codes=[{2}])" -f $case.Name, $result.decision, ($result.reason_codes -join ','))
}

if ($failures -gt 0) {
    Write-Output ("`n{0} case(s) failed." -f $failures)
    exit 1
} else {
    Write-Output "`nAll cases passed."
    exit 0
}
