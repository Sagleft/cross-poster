
![screenshot](screen.png)

Download: https://github.com/Sagleft/cross-poster/releases

Forum topic: https://talk.u.is/viewtopic.php?pid=4627

### How to start

```bash
cp config.example.json config.json
```

update config in `config.json`:

* `utopia.token` - your API token from Utopia API Settings;
* `utopia.channel_id` - the ID of the channel in which the account is an administrator or moderator with the right to publish messages;
* `telegram.token` - your Telegram Bot token given from [@BotFather](https://t.me/BotFather);
* `telegram.chat_id` - your channel (chat) ID to post messages. You can find it out by using [@getmyid_bot](https://t.me/getmyid_bot);

### Build & run

```go
go build
./tool
```

App will be available at: `127.0.0.1:8080/` by default

---
[![udocs](https://github.com/Sagleft/ures/blob/master/udocs-btn.png?raw=true)](https://udocs.gitbook.io/utopia-api/)
