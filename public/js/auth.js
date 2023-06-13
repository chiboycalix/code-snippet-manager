function handleLogout(){
  fetch('/auth/logout', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    }
  }).then((res) => {
    if(res.status === 200){
      window.location.replace('/login')
    } else {
      alert('Could not log out')
    }
  })
}