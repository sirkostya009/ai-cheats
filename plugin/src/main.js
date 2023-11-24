try {
  document.addEventListener('visibilitychange', () => {});

  const decorate = (content) => {
    const gptResponse = document.createElement('h2');
    gptResponse.innerText = content;

    const form = document.getElementById('answ');

    form.parentNode.insertBefore(gptResponse, form.nextSibling);
  };

  const [b3] = document.getElementsByClassName('b3');

  const [question, options] = b3.children;

  fetch('https://ai-cheats-be-dev-mzge.1.ie-1.fl0.io/1', {
    method: 'POST',
    mode: 'cors',
    body: question.innerText +
      [...options.children].reduce((result, { innerText }, i) => `${result}\n${i + 1} ${innerText}`, '')
  }).then(r => r.text())
    .then(decorate)
    .catch(() => {});
} catch (e) {
}
