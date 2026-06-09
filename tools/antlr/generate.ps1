$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$jar = Join-Path $PSScriptRoot "antlr-4.13.2-complete.jar"
$grammarDir = Join-Path $repoRoot "internal\query\parser\antlr\grammar"
$outputDir = Join-Path $repoRoot "internal\query\parser\antlr\generated"

if (-not (Test-Path $jar)) {
    throw "Missing ANTLR jar at $jar"
}

$java = $env:JAVA_HOME
if ($java) {
    $java = Join-Path $java "bin\java.exe"
}

if (-not $java -or -not (Test-Path $java)) {
    $candidates = @(
        "C:\Program Files\Microsoft\jdk-17.0.18.8-hotspot\bin\java.exe",
        "C:\Program Files\Eclipse Adoptium\jdk-17.0.18.8-hotspot\bin\java.exe",
        "C:\Program Files\Android\openjdk\jdk-21.0.8\bin\java.exe",
        "java"
    )

    $java = $candidates | Where-Object {
        if ($_ -eq "java") { return $true }
        Test-Path $_
    } | Select-Object -First 1
}

Push-Location $grammarDir
try {
    & $java -jar $jar -Dlanguage=Go -visitor -package generated -o "..\generated" OutboundAPILexer.g4 OutboundAPIParser.g4
}
finally {
    Pop-Location
}
