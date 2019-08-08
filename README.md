# `oima`

`oima` (**O**CI/ Docker **I**mage Signature **Ma**nagemet Tool/ CLI) is a CLI that helps to manage OCI/ Docker signatures.

## Motivation

We have our signatures in two places: on a _Notary-Server_ and an _S3 Bucket_.

We use the _S3 Bucket_ because of the Pull Signature Check functionality of _CRI-O_.

So it's a huge effort to manage all signatures distributed to two places.
For example, if we update one of our images, then the old image shouldn't be executed anymore in our K8s-Cluster.
So then we have to delete the signatures of the old image from the S3 Bucket _and_ from the Notary-Server.
Also, the signatures of the images are saved with the content digest from Docker in this Format:
`[IMAGE_NAME]@[HASH_ALGO]=[CONTENT_DIGEST]` for example: `hello-world@sha256=92c7f9c92844bbbb5d0a101b22f7c2a7949e40f8ea90c8b3bc396879d95e899a`.


## Usage

This CLI does not have any sub-commands (coming soon), but it has a working terminal UI.

```bash
oima Manages OCI/ Docker Image Signatures in your 'sigstore'.

It's impossible to keep track of all signatures.

Example: you have to remove the signature for the
Docker image 'docker.io/library/hello_world:vulnerable'
- then you have to determine the digest of the image and
manually delete the directory/ dignature.

This tool automates this process and helps to keep
track of all signed images.

Usage:
  oima <command> [flags]
  oima [command]

Available Commands:
  conf        Get configuration variables.
  help        Help for any command.
  image       Interact with images of the remote registry.

Flags:
      --config string   Which config file to use (default is $HOME/.oima.yaml).
      --debug           Print debug messages (defaults to false).
  -h, --help            Display help for oima.
      --version         Display ersion of oima.

Use "oima [command] --help" for more information about a command.
```

To get started, download a release and create a configuration file in `$HOME/.oima.yaml`.
A sample configuration is located in [`examples/`](examples/oima.yaml).
The configuration file is self-explanatory.

Now run the application without any arguments (`oima`), you should now see a "UI".

Keyboard Strokes:
```
q, Ctrl+C               Quit. Exit the application.
e, E                    Exit the image info UI (only works in the image info UI).
d, D                    Delete the signature of a tag (only works in the image info UI).
i, I                    Open the image info UI.
Enter, Space            Expand/ collapse a tree node.
<Arrow Keys>            Move up/ down in the tree or the image info UI.
```


### `Image Info UI`

All tags of an image are listed in the _Image Info UI_.
Here you can check if a tag is signed (or has a Signature) and delete signatures.