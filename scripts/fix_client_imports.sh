#!/bin/bash


sed --version > /dev/null || (echo "Seems like you don’t have gnu-sed" && exit 1)
for f in *; do sed -i '/validators_validator_pb/d' "$f"; done
for f in *; do sed -i '/validate_pb/d' "$f"; done
for f in *; do sed -i '/annotations_pb/d' "$f"; done
for f in *; do sed -i '/scope_pb/d' "$f"; done
