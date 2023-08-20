# Copyright 2023 chenmingyong0423

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SOURCE_COMMIT=../.github/pre-commit
TARGET_COMMIT=../.git/hooks/pre-commit

# copy pre-commit file if not exist.
echo "setting git pre-commit hooks..."
cp $SOURCE_COMMIT $TARGET_COMMIT

# add permission to TARGET_PUSH and TARGET_COMMIT file.
test -x $TARGET_COMMIT || chmod +x $TARGET_COMMIT

echo "installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

echo "installing goimports..."
go install golang.org/x/tools/cmd/goimports@latest
