# tgsend
tgsend is a very basic Telegram client which just sends messages and/or files.
It's mostly useful for sending stuff via CLI or in simple scripts.
However, keep in mind that scripting with clients is against Telegram TOS and should be limited at very specific situations.
More info on this in the [library documentation](https://github.com/gotd/td/blob/main/.github/SUPPORT.md#how-to-not-get-banned).



## Installation
Simply download the [latest release](https://github.com/sgorblex/tgsend/releases) or compile from source:
```sh
go install -v github.com/sgorblex/tgsend
```



## Usage
Initialize the client with a ID/hash API pair (obtainable [here](https://my.telegram.org/)):
```sh
tgsend -init
```
Send a message
```sh
tgsend @person "nuntium mirabilem sane detexi, hanc paginae exiguitas non caperet"
```
Also refer to `tgsend -h`.

The first time you try to send a message, the client will ask you your Telegram login information. You may delete the session file or the client credentials from `$XDG_DATA_HOME/tgsend` (e.g. `~/.local/share/tgsend`).



## Contribute
External contributions are welcome via [pull request](https://github.com/sgorblex/tgsend/pulls). If you are willing to contribute, [TODO.md](TODO.md) is a good place to start.



## License
[MIT](LICENSE)
