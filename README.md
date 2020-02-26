# Go and Vault demo application

This example shows how to use [Okteto](https://github.com/okteto/okteto) and [Okteto Cloud](https://cloud.okteto.com) to develop a Go services that consumes secrets from [Vault](https://www.vaultproject.io/) directly in Kubernetes.

## Prerequisites
- An Okteto Cloud account (sign up for free at https://cloud.okteto.com)
- [Okteto CLI](https://github.com/okteto/okteto) installed locally
- Basic understanding of Vault

## Configure your Okteto Cloud Credentials
If you haven't yet, log in to `https://cloud.okteto.com` and click on the `Credentials` button on the left to download your Okteto Cloud credentials, and set them as your current context by running the command below:

```console
$ export KUBECONFIG=<download-path>/okteto-kube.config
```

## Install Vault

Clone the official vault chart

```console
$ git clone https://github.com/hashicorp/vault-helm
```

Generate the installation configuration file

```console
$ cat <<EOF > config.yaml
injector:
  enabled: false

server:
  dev:
    enabled: true
  authDelegator:
    enabled: false

ui:
  enabled: true
  annotations:
    dev.okteto.com/auto-ingress: "true"
EOF
```

Install the Vault Helm chart

```console
$ helm install vault vault-helm -f config.yaml
```

> For the purpose of this guide, we are installing Vault in `dev` mode. This is not recommended for production scenarios.

## Configure Vault and Store the Secret

Run the vault client

```console
$ kubectl exec -ti vault-0 /bin/sh
```

Login with the root token

```console
$ vault login <root token>
```

Create the policy

```console
cat <<EOF > /home/vault/app-policy.hcl
path "secret*" {
  capabilities = ["read"]
}
EOF

$ vault policy write app /home/vault/app-policy.hcl
```

Create the token

```console
$ vault token create -policy=app
```

Create the secret

```console
$ vault kv put secret/helloworld username=foobaruser password=foobarbazpass
```

## Deploy the application
```console
$ kubectl apply -f k8s.yaml
```

The default application uses hard-coded credentials instead of getting them from Vault. Get the URL from the Okteto Cloud UI, and call the application from your local terminal:

```console
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
```console
$ kubectl create secret generic vault --from-literal=token=<vault token>
```

And update your application:
```console
$ kubectl apply -f k8s-with-secret.yaml
```

## Develop directly in Kubernetes

Now let's go ahead and integrate our code with Vault, directly in the cluster.

Launch your development environment

```console
$ okteto up
```

When you launch a development environment in Okteto, it will inherit all the configurations of the original deployment, like the Vault API token we just added:

```console
okteto> echo $VAULT_TOKEN
```

You'll notice that the code to get the secrets is already there. All you need to do to enable is to open `main.go` in your local IDE and update the `helloServer` function to call `getVaultSecrets()` instead.

```golang
func helloServer(w http.ResponseWriter, r *http.Request) {
	u, p, err := getHardcodedSecrets()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		w.WriteHeader(500)
	}

	fmt.Fprint(w, "Hello world!\n")
	fmt.Fprintf(w, "Super secret username: %s\n", u)
	fmt.Fprintf(w, "Super secret password is: %s\n", p)
}
```

Go to your okteto terminal, and start the service:
```console
okteto> go run main.go
```

```console
Starting hello-world server..
```

Call the application again to see the changes:

```console
$ curl https://hello-world-cindy.cloud.okteto.net
```

```console
Hello world!
Super secret username: foobaruser
Super secret password is: foobarbazpass
```

# Conclusion

In this guide we show you how you can integrate your application with a third-party service like Vault using Okteto and Okteto Cloud. 

Being able to use the same configuration model as you have in production is one of the great advantages of using remote development environments. It helps you go faster, better understand your dependencies, and catch pesky integration bugs early in the cycle. 

Visit https://okteto.com/docs to learn how Okteto and Okteto Cloud makes application development easier than ever.
