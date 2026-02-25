param(
    [string]$Tag = 'v0.1.0',
    [string]$TargetCommitish
)

if ($env:GITHUB_TOKEN -eq $null -or $env:GITHUB_TOKEN -match 'PASTE_TOKEN_HERE') {
    Write-Error 'GITHUB_TOKEN not set or placeholder. Set $env:GITHUB_TOKEN and re-run.'
    exit 2
}

$changelogPath = Join-Path $PSScriptRoot '..\CHANGELOG.md'
if (!(Test-Path $changelogPath)) {
    Write-Error "CHANGELOG.md not found at $changelogPath"
    exit 2
}

$changelog = Get-Content $changelogPath -Raw
$escapedTag = [regex]::Escape($Tag)
$pattern = "(?ms)^##\s*\[?$escapedTag\]?.*?(?=^##\s|\z)"
$match = [regex]::Match($changelog, $pattern)
if (-not $match.Success) {
    Write-Error "$Tag section not found in CHANGELOG.md"
    exit 2
}

$notes = $match.Value.Trim()
$notesFileName = "RELEASE_NOTES_${($Tag -replace '[^A-Za-z0-9_.-]','_')}.md"
$notesPath = Join-Path $PSScriptROOT "..\$notesFileName"
Set-Content -Path $notesPath -Value $notes -Encoding UTF8
Write-Output "Wrote release notes to $notesPath"

 $owner = 'BrunoBaumgartner78'
 $repo = 'digiemu-core'
 # $Tag param is used

$headers = @{ Authorization = ('Bearer ' + $env:GITHUB_TOKEN); 'X-GitHub-Api-Version' = '2022-11-28'; Accept = 'application/vnd.github+json'; 'User-Agent' = 'digiemu-core-release' }

try {
    # Check if a release already exists for this tag
    $existing = $null
    try {
        $existing = Invoke-RestMethod -Headers $headers -Uri "https://api.github.com/repos/$owner/$repo/releases/tags/$Tag"
    } catch {
        $existing = $null
    }

    $payload = @{ tag_name = $Tag; name = $Tag; body = $notes; draft = $true; prerelease = $false }
    if ($TargetCommitish) { $payload.target_commitish = $TargetCommitish }
    $body = $payload | ConvertTo-Json -Depth 6

    if ($existing -ne $null) {
        Write-Output "Existing release found id=$($existing.id) url=$($existing.html_url) - will PATCH to update body and ensure draft=true"
        $rid = $existing.id
        $uri = "https://api.github.com/repos/$owner/$repo/releases/$rid"
        $r = Invoke-RestMethod -Method Patch -Uri $uri -Headers $headers -Body $body
    } else {
        $uri = "https://api.github.com/repos/$owner/$repo/releases"
        $r = Invoke-RestMethod -Method Post -Uri $uri -Headers $headers -Body $body
    }

    Write-Output ("DRAFT_RELEASE_ID=" + $r.id)
    Write-Output ("DRAFT_RELEASE_HTML_URL=" + $r.html_url)
    Write-Output ("DRAFT_RELEASE_EDIT_URL=https://github.com/$owner/$repo/releases/edit/" + $r.id)
    Write-Output ("DRAFT_RELEASE_TAG_URL=https://github.com/$owner/$repo/releases/tag/$Tag")
    exit 0
} catch {
    Write-Error 'Release creation/update failed.'
    if ($_.Exception.Response -ne $null) {
        try {
            $sr = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $err = $sr.ReadToEnd()
            Write-Error $err
        } catch {}
    }
    exit 3
} finally {
    if (Test-Path $notesPath) { Remove-Item -Force $notesPath }
}
