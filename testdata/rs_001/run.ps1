param(
    [string]$TestDir = ""
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($TestDir)) {
    $TestDir = $PSScriptRoot
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")

$schemaIntent  = Join-Path $repoRoot "schemas\intent-envelope.schema.json"
$schemaAdm     = Join-Path $repoRoot 'schemas/admission_result_v0.2.schema.json'
$regRules      = Join-Path $repoRoot 'admission-rule-registry.yaml'
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

function Sha256($s) {
    $b = [System.Text.Encoding]::UTF8.GetBytes($s)
    return [BitConverter]::ToString([System.Security.Cryptography.SHA256]::Create().ComputeHash($b)).Replace('-', '').ToLower()
}

function Canonical-Value($v) {
    if ($v -eq $null) { return $null }
    if ($v -is [array] -or ($v -is [System.Collections.IEnumerable] -and -not ($v -is [string]))) {
        return @($v | ForEach-Object { Canonical-Value $_ })
    }
    if ($v -is [System.Collections.Specialized.OrderedDictionary] -or $v -is [hashtable]) {
        $o = [ordered]@{}; foreach ($k in ($v.Keys | Sort-Object)) { $o[$k] = Canonical-Value $v[$k] }; return $o
    }
    if ($v.PSObject -and $v.PSObject.Properties -and -not ($v -is [string] -or $v -is [int] -or $v -is [double] -or $v -is [bool])) {
        $o = [ordered]@{}; foreach ($k in ($v.PSObject.Properties.Name | Sort-Object)) { $o[$k] = Canonical-Value $v.$k }; return $o
    }
    return $v
}

function Get-IntentDigest($intent) {
    $o = [ordered]@{
        intent_digest_profile = 'P0.ADMISSION.INTENT.v0.1'
        schema_version        = $intent.schema_version
        architecture_revision = $intent.architecture_revision
        capability_ref        = $intent.capability_ref
        aggregate_ref         = $intent.aggregate_ref
        command_ref           = $intent.command_ref
        payload               = (Canonical-Value $intent.payload)
    }
    $json = ConvertTo-Json -InputObject $o -Compress -Depth 10
    return 'p0-intent:sha256:' + (Sha256 $json)
}

function Get-AdmissionId($arch, $intentDigest, $capRef, $aggRef, $cmdRef, $transRef, $decision, $rules, $reasons) {
    $o = [ordered]@{
        admission_id_profile = 'P0.ADMISSION.ID.v0.1'
        schema_version       = 'v0.1'
        architecture_revision = $arch
        intent_digest        = $intentDigest
        capability_ref       = $capRef
        aggregate_ref        = $aggRef
        command_ref          = $cmdRef
        transition_ref       = $transRef
        decision             = $decision
        rule_refs            = @($rules | Sort-Object)
        reason_codes         = @($reasons | Sort-Object)
    }
    $json = ConvertTo-Json -InputObject $o -Compress -Depth 10
    return 'admission:sha256:' + (Sha256 $json)
}

function New-AdmissionResult($admissionId, $decision, $capRef, $aggRef, $cmdRef, $transRef, $reasons, $rules) {
    $o = [ordered]@{
        schema_version        = 'v0.2'
        architecture_revision = $baselineRevision
        admission_id          = $admissionId
        decision              = $decision
        capability_ref        = $capRef
        aggregate_ref         = $aggRef
        command_ref           = $cmdRef
        rule_refs             = @($rules | Sort-Object)
        reason_codes          = @($reasons | Sort-Object)
    }
    if ($decision -eq 'ADMIT') { $o.transition_ref = $transRef }
    return $o
}

function Parse-Rules($path) {
    $rules = @(); $current = $null; $inRules = $false; $inMap = $false
    foreach ($line in (Get-Content -Path $path)) {
        $line = $line.TrimEnd()
        if ($line -eq 'admission_rules:') { $inRules = $true; $inMap = $false; continue }
        if ($line -eq 'reason_code_to_rule:') { $inRules = $false; $inMap = $true; continue }
        if ($line -eq 'rejected_reason_codes:') { $inRules = $false; $inMap = $false; continue }
        if ($inRules -and $line.StartsWith('  - rule_id: ')) { $current = $line.Substring(13).Trim(); $rules += @{ rule_id = $current } }
        if ($inRules -and $line.StartsWith('    failure_reason_code: ')) { $rules[$rules.Count - 1].failure_reason_code = $line.Substring(25).Trim() }
    }
    return $rules
}

function Parse-ReasonToRule($path) {
    $reasonToRule = @{}; $current = $null; $inMap = $false
    foreach ($line in (Get-Content -Path $path)) {
        $line = $line.TrimEnd()
        if ($line -eq 'reason_code_to_rule:') { $inMap = $true; continue }
        if ($inMap -and $line.StartsWith('  ') -and $line.Trim().EndsWith(':') -and -not $line.StartsWith('    -')) { $current = $line.Trim().TrimEnd(':'); continue }
        if ($inMap -and $line.StartsWith('    - ')) { $reasonToRule[$current] = $line.Substring(6).Trim() }
    }
    return $reasonToRule
}

function Get-PropertyOrDefault($obj, $name, $default) {
    if ($obj -and $obj.PSObject.Properties[$name]) { return $obj.$name }
    return $default
}

$capabilities     = Parse-Capabilities $regCapability
$ownership        = Parse-Ownership $regAggregate
$commands         = Parse-Commands $regCommand
$baselineRevision = Get-BaselineRevision $regBaseline
$rules            = Parse-Rules $regRules
$reasonToRule     = Parse-ReasonToRule $regRules

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
    $ruleIds = @()
    $decision = 'ADMIT'
    $capRef   = Get-PropertyOrDefault $intent 'capability_ref' ''
    $aggRef   = Get-PropertyOrDefault $intent 'aggregate_ref'  ''
    $cmdRef   = Get-PropertyOrDefault $intent 'command_ref'    ''
    $transRef = $null
    $admissionId = $null
    $intentDigest = $null

    if (-not $validIntent) {
        $ruleIds += 'P0.ADMISSION.INTENT_REQUIRED_FIELDS'
        if ($intent -and -not $intentDigest) { $intentDigest = Get-IntentDigest $intent }
        if ($intentErr -match 'Required properties') {
            $reasons += 'MISSING_REQUIRED_FIELD'
        } else {
            $reasons += 'INVALID_INTENT'
        }
        $decision = 'REJECT'
    } else {
        $intentDigest = Get-IntentDigest $intent
        $failed = $false
        foreach ($rule in $rules) {
            $reason = $null
            switch ($rule.rule_id) {
                'P0.ADMISSION.ARCHITECTURE_REVISION' {
                    if ($intent.architecture_revision -ne $baselineRevision) { $reason = $rule.failure_reason_code }
                }
                'P0.ADMISSION.CAPABILITY_EXISTS' {
                    if (-not $capabilities.ContainsKey($capRef)) { $reason = $rule.failure_reason_code }
                }
                'P0.ADMISSION.CAPABILITY_MUTATES' {
                    if (-not $capabilities[$capRef]) { $reason = $rule.failure_reason_code }
                }
                'P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY' {
                    $owner = $null
                    foreach ($agg in $ownership.Keys) { if ($ownership[$agg] -contains $capRef) { $owner = $agg; break } }
                    if ($owner -ne $aggRef) { $reason = $rule.failure_reason_code }
                }
                'P0.ADMISSION.COMMAND_EXISTS' {
                    if (-not $commands.ContainsKey($cmdRef)) { $reason = $rule.failure_reason_code }
                }
                'P0.ADMISSION.COMMAND_CAPABILITY_MATCH' {
                    if ($commands.ContainsKey($cmdRef) -and $commands[$cmdRef].capability_id -ne $capRef) { $reason = $rule.failure_reason_code }
                }
                'P0.ADMISSION.COMMAND_AGGREGATE_MATCH' {
                    if ($commands.ContainsKey($cmdRef) -and $commands[$cmdRef].aggregate_id -ne $aggRef) { $reason = $rule.failure_reason_code }
                }
                'P0.ADMISSION.COMMAND_TRANSITION_DEFINED' {
                    if ($commands.ContainsKey($cmdRef) -and [string]::IsNullOrWhiteSpace($commands[$cmdRef].transition_id)) {
                        $reason = $rule.failure_reason_code
                    } elseif ($commands.ContainsKey($cmdRef)) {
                        $transRef = $commands[$cmdRef].transition_id
                    }
                }
                'P0.ADMISSION.INTENT_REQUIRED_FIELDS' {
                }
            }
            $ruleIds += $rule.rule_id
            if ($reason) { $reasons += $reason; $failed = $true; break }
        }
        if ($failed) { $decision = 'REJECT' }
    }

    if ($intent -and -not $intentDigest) { $intentDigest = Get-IntentDigest $intent }
    $admissionId = Get-AdmissionId $intent.architecture_revision $intentDigest $capRef $aggRef $cmdRef $transRef $decision $ruleIds $reasons

    $result = New-AdmissionResult $admissionId $decision $capRef $aggRef $cmdRef $transRef $reasons $ruleIds
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

    if ($decision -eq 'ADMIT') {
        $cmdEnv = [ordered]@{
            schema_version        = 'v0.1'
            architecture_revision = $baselineRevision
            command_id            = $admissionId
            admission_ref         = $admissionId
            capability_ref        = $capRef
            aggregate_ref         = $aggRef
            command_ref           = $cmdRef
            transition_ref        = $transRef
            payload               = (Canonical-Value $intent.payload)
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

    Write-Output ('{0}: PASS (decision={1}, reason_codes=[{2}], admission_id={3})' -f $case.Name, $result.decision, ($result.reason_codes -join ','), $admissionId)
}

$schemaTestDir = Join-Path $TestDir '_schema_validation'
if (Test-Path $schemaTestDir) {
    $admit = Get-Content (Join-Path $schemaTestDir 'admit_no_transition.json') -Raw
    $reject = Get-Content (Join-Path $schemaTestDir 'reject_no_transition.json') -Raw
    $admitValid = $false
    try { $admitValid = Test-Json -Json $admit -SchemaFile $schemaAdm } catch { $admitValid = $false }
    if ($admitValid) {
        Write-Output 'SCHEMA admit_no_transition: FAIL (expected invalid)'
        $failures++
    } else {
        Write-Output 'SCHEMA admit_no_transition: PASS'
    }
    $rejectValid = $false
    try { $rejectValid = Test-Json -Json $reject -SchemaFile $schemaAdm } catch { $rejectValid = $false }
    if ($rejectValid) {
        Write-Output 'SCHEMA reject_no_transition: PASS'
    } else {
        Write-Output 'SCHEMA reject_no_transition: FAIL'
        $failures++
    }
}

if ($failures -gt 0) {
    Write-Output ('{0} case(s) failed.' -f $failures)
    exit 1
} else {
    Write-Output 'All cases passed.'
    exit 0
}
