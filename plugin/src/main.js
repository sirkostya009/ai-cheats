(async () => {
  try {
    const [b3] = document.getElementsByClassName('b3');

    const [question, options] = b3.children;
    const body = [...(question.innerText +
      [...options.children].reduce((result, {innerText}, i) => `${result}\n${i + 1} ${innerText}`, ''))];

    const response = fetch('https://ai-cheats.2.ie-1.fl0.io/1', {
      method: 'POST',
      mode: 'cors',
      body: body.reduce((result, char, i) =>
        `${result}${body[i]}${require('object-hash').MD5(document.getElementsByClassName('b2')[0].innerText)[i] || '_'}`, '')
    });

    const answ = document.getElementById('answ');
    const answer = document.createElement('h2');
    answ.parentNode.insertBefore(answer, answ.nextSibling);
    answer.innerText = await (await response).text();
  } catch (e) {
  }
})();
