function handleLogout(){
  fetch('/auth/logout', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    }
  }).then((res) => {
    if(res.status === 200){
      localStorage.removeItem('codeSnippetManagerJWT')
      window.location.replace('/login')
    } else {
      alert('Could not log out')
    }
  })
}

function onLogin(event){
  console.log('onLogin')
  fetch('/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: document.getElementById('email').value,
      password: document.getElementById('password').value
    })
  }).then((res) => {
    return res.json()
  }).then((data) => {
    if(data.error){
      alert(data.error)
    } else {
      localStorage.setItem('codeSnippetManagerJWT', data.token)
      window.location.replace('/')
    }
  }).catch((err) => {
    console.log(err)
  })
}