# `oima`

`oima` (**O**CI/ Docker **I**mage Signature **Ma**nagemet Tool/ CLI) is a CLI that helps managing OCI/ Docker Signatures.

## Motivation

We have at two places our Signatures: on a _Notary-Server_ and an _S3 Bucket_.

We use the _S3 Bucket_ because of the Pull Signature Check functionality of _CRI-O_.

So it's a huge expenditure to manage all Signatures which are distributed at two places.
For example, if we update one of our Images, then the old Image shouldn't be executed anymore in our K8s-Cluster.
Now we have to delete the Signatures of the old Image from the S3 Bucket _and_ from the Notary-Server.
And the signatures of the images are saved with the Content Digest from Docker in this Format:
`[IMAGE_NAME]@[HASH_ALGO]=[CONTENT_DIGEST]` for example: `hello-world@sha256=92c7f9c92844bbbb5d0a101b22f7c2a7949e40f8ea90c8b3bc396879d95e899a`.


## Usage

This CLI does not have any sub-commands (coming soon), but it has a working Terminal UI.

```bash
oima Manages OCI/ Docker Image Signatures in you 'sigstore'.

It's impossible to keep track of all Signatures.

For Example, you have to remove the Signature for the
Docker Image 'docker.io/library/hello_world:vulnerable',
now you have to find out the Digest of the Image and
manually delete the Directory/ Signature.

This Tool automates this Process and helps to keep
track of all signed Images.

Usage:
  oima <command> [flags]
  oima [command]

Available Commands:
  conf        Get Configuration Variables
  help        Help about any command
  image       Interact with Images of the Remote Registry

Flags:
      --config string   config file (default is $HOME/.oima.yaml)
      --debug           Print Debug Messages (defaults to false)
  -h, --help            help for oima
      --version         version for oima

Use "oima [command] --help" for more information about a command.
```

To get started download a release and create a Config File at `$HOME/.oima.yaml`.
A Config Example is located under [`examples/`](examples/oima.yaml).
The Config File is self-explaining.

Now run the Application without any Arguments (`oima`), you should now see a "UI".

Keyboard Strokes:
```
q, Ctrl+C               Quit. Exit the Application
e, E                    Exit the Image Info UI (only works in the Image Info UI)
d, D                    Delete a Signature of a Tag (only works in the Image Info UI)
i, I                    Open the Image Info UI
Enter, Space            Expand/ Collapse a Tree Node
<Arrow Keys>            Move Up/ Down in the Tree or the Image Info UI
```


### `Image Info UI`

In the _Image Info UI_ are all Tags of an Image listed.
Here you can check if a Tag is signed (or has a Signature) and delete Signatures.