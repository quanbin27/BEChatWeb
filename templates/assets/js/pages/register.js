function sendRegisterInfo(info){
    console.log(info)
    fetch('/api/v1/register',{
        method:'POST',
        headers:{
            'content-type':'application/json'
        },
        body:JSON.stringify(info)
    }).then(async response=>{
        if(!response.ok){
            const notification = document.querySelector('.notification');
            notification.textContent = 'Không thể đăng ký vui lòng thử lại';
            notification.style.display = '';
        }
        const responseJson = await response.json();
        if(responseJson.Status.toLowerCase() == 'success'){
            window.alert('Register successfully');
            window.location.href = '/login';
        }
    })
}
function main(){
    const emailIput = document.querySelector('#useremail')
    const usernameInput = document.querySelector('#username')
    const passwordInput = document.querySelector('#userpassword')
    const registerForm = document.querySelector('#register-form');
    const btnSubmit = registerForm.querySelector('button');
    btnSubmit.onclick = (e)=>{
        e.preventDefault();
        if(usernameInput.value === ''||emailIput.value === '' || passwordInput.value === '')
            return;
        const info = {
            name:usernameInput.value,
            email: emailIput.value,
            password: passwordInput.value
        }
        sendRegisterInfo(info);
    }
}
main();