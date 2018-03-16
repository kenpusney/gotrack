#!/bin/sh
set -e

wget -nc https://github.com/kenpusney/rescenario/releases/download/v0.0.3/ReScenario-0.0.3.jar

java -jar ReScenario-0.0.3.jar scenarios/TestKana.yml