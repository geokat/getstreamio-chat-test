async function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function runLoopForUser(uid) {
  let resp = await fetch(`http://localhost:8000/users/${uid}/token`);
  const token = await resp.text();

  const client = StreamChat.getInstance('{{.APIKey}}');

  await client.connectUser(
    {
      id: uid,
      name: uid
    },
    token
  );

  const channel = client.channel('messaging', 'talk', {
    name: 'Test channel',
  });

  await channel.watch();

  await channel.on('message.new', e => {
    console.log(`${uid} got message:`, e.message.text);
  });

  for (let i = 0; i < 100; i++) {
    await sleep(1000);
    await channel.sendMessage({
      text: `message ${i} from ${uid}`
    });
  }
}
