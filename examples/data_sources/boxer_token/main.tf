terraform {
  required_providers {
    boxer = {
      source = "registry.terraform.io/sneaksAndData/boxer"
    }
  }
}

provider "boxer" {
  issuer_host = "http://localhost:8888/"
}

resource "boxer_cedar_schema" "example"  {
  id = "example"
  data_json = <<EOT
  {
    "PhotoApp": {
      "commonTypes": {
        "PersonType": {
          "type": "Record",
          "attributes": {
            "age": {
              "type": "Long"
            },
            "name": {
              "type": "String"
            }
          }
        }
      },
      "entityTypes": {
        "User": {
          "memberOfTypes": [
            "UserGroup"
          ],
          "shape": {
            "type": "Record",
            "attributes": {
              "personInformation": {
                "type": "PersonType"
              },
              "userId": {
                "type": "String"
              }
            }
          }
        },
        "UserGroup": {}
      },
      "actions": {}
    }
  }
EOT
}

resource "boxer_principal" "example" {
  schema_id = boxer_cedar_schema.example.id
  data_json = <<EOT
{
    "uid": {
        "type": "PhotoApp::User",
        "id": "alice"
    },
    "attrs": {
        "userId": "897345789237492878",
        "personInformation": {
            "age": 85,
            "name": "alice"
        }
    },
    "parents": [
        {
            "type": "PhotoApp::UserGroup",
            "id": "alice_friends"
        },
        {
            "type": "PhotoApp::UserGroup",
            "id": "AVTeam"
        }
    ]
}
EOT
}

resource "boxer_identity_provider" "example"  {
  name = "provider"
  user_id_claim = "preferred_username"
  discovery_url = "http://localhost:8080/realms/master/"
  issuers = [
    "http://localhost:8080/realms/master",
  ]
  audiences = [
    "account"
  ]
}

resource "boxer_external_identity" "example" {
  identity_provider = "provider"
  id                   = "test_user"
  principal = {
    schema_id = boxer_cedar_schema.example.id
    principal_id = boxer_principal.example.id
  }
}

data boxer_token "example" {
  identity_provider = boxer_identity_provider.example.name
  auth = {
    header = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJPY0x4RTJfdGVxZGdSNzNyN2UtZHBwSUtBUXJManZXNXZ6NWxVbWVGQjVvIn0.eyJleHAiOjE3NTM5Njk5NTYsImlhdCI6MTc1Mzk2OTg5NiwianRpIjoib25ydHJvOjMzYWY3NWZlLTExNTgtNGEyMy1iMjZmLThhMDkyYjYxZjJlMiIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODA4MC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiYWNjb3VudCIsInR5cCI6IkJlYXJlciIsImF6cCI6InRlc3RfY2xpZW50Iiwic2lkIjoiYWVkZTYyY2MtYWJhMS00Nzk5LThkNjItNmY4M2M3ZTY3MmViIiwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImRlZmF1bHQtcm9sZXMtbWFzdGVyIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdF91c2VyIn0.CMrLLBd3Au9trXUQN1lUwA-YcWMxqwMKpIsBp7-PAkE2iEgeGLFxQh2mojKIy7OrlROSJomuAPU9g03oIx5vBPPp3NbAYjt1CECqCFGJVnyPn2qufOrYD0EOpYRgYUPswwcuO-nN5SWeHEVTX2nj7KzpIDB5Frf2F4r-Sj35oXgBQYFYpsKd57o6kg-HzHigSDuSDVrcrjqs-eUUBuDyp1vkdhF8n8gZ1VLFbU8wO5WwhLhwcBTYp7Bni71UNwjM4vJx41oqnW0YKtTPR8tahKDZVCTrZqTTuRwxVgBiw-Aa52O2nSV0Vivuh_CnKMkgLxPoqasFe6A0Yikrr2z5yA"
  }
}

output "test" {
  value = data.boxer_token.example
}
