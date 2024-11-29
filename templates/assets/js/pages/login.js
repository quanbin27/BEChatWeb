function login(email,password){
    fetch('/api/v1/login',{
        method:'POST',
        headers:{
            'content-type':'application/json',
        },
        body:JSON.stringify({
            "email":email,
            "password":password,
        })
    }).then(async response=>{
        if(response.ok){
            const loginInfo = await response.json();
            const token = loginInfo.Token;
            localStorage.setItem('token',token);
            window.location.href = '/';
        }
        else{
            toastr.error('Email or password is incorrect','Fail')
        }
    })
}

function main(){
    toastr.options = {
        "closeButton": false, // Không hiện nút đóng
        "progressBar": true,  // Thanh tiến trình hiển thị thời gian tắt
        "positionClass": "toast-top-right", // Vị trí: góc phải trên
        "timeOut": "2000",    // Tự động đóng sau 3 giây
    };
    const btnLogIn = document.querySelector('.btn.btn-primary.w-100');
    const emailInput = document.querySelector('#email');
    const passwordInput = document.querySelector('#password-input');
    btnLogIn.onclick = (e)=>{
        e.preventDefault();
        if(emailInput.value === ''){
            emailInput.placehoder = 'please fill email';
            toastr.error('Please enter email','Notification');
            return;
        }
        if(passwordInput.value === ''){
            toastr.error('Please enter password','Notification');
            return;
        }
        login(emailInput.value,passwordInput.value);
    }

}
main();