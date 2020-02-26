# Go and Vault demo application

This example shows how to use [Okteto](https://github.com/okteto/okteto) and [Okteto Cloud](https://cloud.okteto.com) to develop a Go services that consumes secrets from Vault.

## Prerequisites
- An Okteto Cloud account (sign up for free at https://cloud.okteto.com)
- [Okteto CLI](https://github.com/okteto/okteto)

## Install Vault

1. Clone the official vault chart
```
$ git clone https://github.com/hashicorp/vault-helm
```

1. Generate the installation configuration file
```
$ cat <<EOF > config.yaml
injector:
  enabled: false

server:
  standalone:
    enabled: true
  authDelegator:
    enabled: false

ui:
  enabled: true
  annotations:
    dev.okteto.com/auto-ingress: "true"
EOF
```

1. Install Vault
```
$ helm install vault vault-helm -f config.yaml
```

## Configure Vault and Store the Secret

1. Run the vault client
```
$ kubectl exec -ti vault-0 /bin/sh
```

1. Login with the root token
```
$ vault login <root token>
```

1. Create the policy
```
cat <<EOF > /home/vault/app-policy.hcl
path "secret*" {
  capabilities = ["read"]
}
EOF

$ vault policy write app /home/vault/app-policy.hcl
```

1. Enable the K/V store
```
$ vault secrets enable -path=secret kv-v2
```

1. Create the token
```
$ vault token create -policy=app
```

1. Create the secret
```
$ vault kv put secret/helloworld username=foobaruser password=foobarbazpass
```

## Deploy the application
```
$ kubectl apply -f k8s.yaml
```

The default application uses hard-coded credentials instead of getting them from Vault. Get the URL from the Okteto Cloud UI, and call the application from your local terminal:

```
$ curl https://hello-world-cindy.cloud.okteto.net
```

```
Hello world!
Super secret username: hard-coded-user
Super secret password is: hard-coded-password
```

## Prepare your Application to use Vault
In order to call the vault API, you'll need to pass a token to your service.  

We'll use the one we created in the previous step, and pass it via an environment variable.

Create the secret:
```
$ kubectl create secret generic vault --from-literal=token=<vault token>
```

And update your application:
```
$ kubectl apply -f k8s-with-secret.yaml
```

## Develop directly in Kubernetes

Now let's go ahead and integrate our code with Vault, directly in the cluster.

Launch your development environment
```
okteto up
```

When you launch a development environment in Okteto, it will inherit all the configurations of the original deployment, like the Vault API token we just added:
```
okteto> echo $VAULT_TOKEN
```

You'll notice that the code to get the secrets is already there. All you need to do to enable is to open `main.go` in your local IDE and update line `55` to call `getVaultSecrets()` instead.

Go to your okteto terminal, and start the service:
```
okteto> go run main.go
```

```
Starting hello-world server..
```

Call the application again to see the changes:

```
$ curl https://hello-world-cindy.cloud.okteto.net
```

```
Hello world!
Super secret username: foobaruser
Super secret password is: foobarbazpass
```

# Conclusion

In this guide we show you how you can integrate your application with a third-party service like Vault using Okteto and Okteto Cloud. 

Being able to use the same configuration model as you have in production is one of the great advantages of using remote development environments. It helps you go faster, better understand your dependencies, and catch pesky integration bugs early in the cycle. 

Visit https://okteto.com/docs to learn how Okteto and Okteto Cloud makes application development easier than ever.
