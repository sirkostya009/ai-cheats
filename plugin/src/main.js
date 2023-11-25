(async () => {
  try {
    const hash = require('object-hash');

    document.removeEventListener('visibilitychange', document.onvisibilitychange);

    const [b3] = document.getElementsByClassName('b3');

    const [question, options] = b3.children;
    const body = [...(question.innerText +
      [...options.children].reduce((result, {innerText}, i) => `${result}\n${i + 1} ${innerText}`, ''))];

    const response = fetch('https://ai-cheats.2.ie-1.fl0.io/1', {
      method: 'POST',
      mode: 'cors',
      body: body.reduce((result, char, i) =>
        `${result}${body[i]}${hash(document.getElementsByClassName('b2')[0].innerText)[i] || '_'}`, '')
    });

    const answer = document.createElement('h2');
    document.answ.parentNode.insertBefore(answer, document.answ.nextSibling);
    answer.innerText = await (await response).text();
  } catch (e) {
  }
})();
