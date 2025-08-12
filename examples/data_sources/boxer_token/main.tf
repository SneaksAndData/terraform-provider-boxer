terraform {
  required_providers {
    boxer = {
      source = "registry.terraform.io/sneaksAndData/boxer"
    }
  }
}

provider "boxer" {
  issuer_host = "http://localhost:8888/"
  validator_host = "http://localhost:8081/"
}

resource "boxer_issuer_cedar_schema" "example"  {
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
  schema_id = boxer_issuer_cedar_schema.example.id
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
    schema_id = boxer_issuer_cedar_schema.example.id
    principal_id = boxer_principal.example.id
  }
}

data boxer_token "example" {
  depends_on = [
    boxer_external_identity.example
  ]
  identity_provider = boxer_identity_provider.example.name
  auth = {
    # For testing purposes, provide the bearer token value manually
    bearer = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI2VG1pRzRfcTNuQS1hOUMycmhkOURoNEdzR1VFSU9SaUI1ZTJFUkJNbDNvIn0.eyJleHAiOjE3NTQ5OTE3MTIsImlhdCI6MTc1NDk5MTY1MiwianRpIjoib25ydHJvOjQzYzg1ZWI2LWIyMjMtNGM0NS1iZDEyLTI4YTUwNmNjNDFiNyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODA4MC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiYWNjb3VudCIsInR5cCI6IkJlYXJlciIsImF6cCI6InRlc3RfY2xpZW50Iiwic2lkIjoiY2ZmMjUyZTItNWJmZS00ZDllLTg3Y2UtMTgzM2QwMjYzMGRlIiwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImRlZmF1bHQtcm9sZXMtbWFzdGVyIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdF91c2VyIn0.RiXVFdz7LHd0-bDd7_yM4gRrr5MHr37ysX0oEr5Eg6vVsioEIFn809fprd7XAodP_jBpjglcD12bhS3JCPxAIx_GQYj0Iqbkunj1PcFL_xkfflQmcc0dRf1izpol6X-K61FvnTtvG5TP_EJspKR8OeOuQCl788ES-vf4WooYuE54wRU4NngqChG7Q7KOrsbNlAjxZMGjaOv9C8fo7oNWNYGH9BSvYrabx06SIJjQjOICXIitO81um70HW3RGA61qMguCF5gkR94_xSrME_dsY1-EWvFQ65I-OTKnLSuh8i_ZwKqpd871m-zyAiFwXuna136nt4fNitEhMTXPpLcIkw"
  }
}

output "test" {
  value = data.boxer_token.example
}
