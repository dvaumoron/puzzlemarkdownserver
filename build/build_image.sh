#!/usr/bin/env bash

go install

buildah from --name puzzlemarkdownserver-working-container scratch
buildah copy puzzlemarkdownserver-working-container $HOME/go/bin/puzzlemarkdownserver /bin/puzzlemarkdownserver
buildah config --env SERVICE_PORT=50051 puzzlemarkdownserver-working-container
buildah config --port 50051 puzzlemarkdownserver-working-container
buildah config --entrypoint '["/bin/puzzlemarkdownserver"]' puzzlemarkdownserver-working-container
buildah commit puzzlemarkdownserver-working-container puzzlemarkdownserver
buildah rm puzzlemarkdownserver-working-container
