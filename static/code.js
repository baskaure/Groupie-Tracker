function Code() {
    var codeElement = document.getElementById('code');
    var tempElement = document.createElement('textarea');
    tempElement.value = codeElement.innerText;
    document.body.appendChild(tempElement);
    tempElement.select();
    document.execCommand('copy');
    document.body.removeChild(tempElement);
    alert('Code Copier !');
}
