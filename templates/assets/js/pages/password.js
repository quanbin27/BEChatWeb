
async function getInfo(token){
        return  fetch('/api/v1/user/me',{
            method:'GET',
            headers:{
                'content-type':'application/json',
                'authorization':`Bearer ${token}`
            }
        }).then(async response=>{
            if(response.ok){
                const userInfo = await response.json()
                return userInfo;
            }
            else{
                //yeu cau dang nhap
                window.location.href = '/login';
            }
        })
}

async function save(token,passwordInfo){
    fetch('/api/v1/changePassword',{
        method:'POST',
        headers:{
            'content-type':'application/json',
            'authorization': `Bearer ${token}`
        },
        body: JSON.stringify(passwordInfo)
    }).then(re=>{
        if(re.ok){
            toastr.success('Change password successfully','Sucessfully');
            setTimeout(()=>{
                window.location.href = '/login';
            },1000);
        }
        else
            toastr.error('Wrong password','Fail');
    })
}
function cancel(){
    window.location.href = '/';
}
async function main(){
    const token = localStorage.getItem('token');
    const name = document.querySelector('.username');
    const avatar = document.querySelector('.rounded-circle.img-thumbnail.avatar-lg');
    const oldPasswordInput = document.querySelector('#oldpassword-input');
    const newPasswordInput = document.querySelector('#password-input');
    const confirmPasswordInput = document.querySelector('#confirmpassword-input');
    const btnSave = document.querySelector('.btn.btn-primary.w-100');
    const btnCancel = document.querySelector('.btn.btn-light.w-100');


    toastr.options = {
        "closeButton": false, // Không hiện nút đóng
        "progressBar": true,  // Thanh tiến trình hiển thị thời gian tắt
        "positionClass": "toast-top-right", // Vị trí: góc phải trên
        "timeOut": "3000",    // Tự động đóng sau 3 giây
    };

    const info = await getInfo(token);

    name.textContent = info.Name;
    if(info.avatar!=null){
        avatar.src = info.avatar;
    }
    btnCancel.onclick = ()=>cancel();
    btnSave.onclick = (e)=>{
        e.preventDefault();
        if(oldPasswordInput.value == ''){
            toastr.error('Please enter old password','Notification')
            return;
        }
        if(newPasswordInput.value == ''){
            toastr.error('Please enter new password','Notification')
            return;
        }
        if(confirmPasswordInput.value == ''){
            toastr.error('Please enter confirm password','Notification')
            return;
        }
        if(confirmPasswordInput.value != newPasswordInput.value){
            toastr.error('Confirm password and new password do not match');
            return;
        }
        save(token,{
            old_password: oldPasswordInput.value,
            new_password: newPasswordInput.value,
            confirm_new_password: confirmPasswordInput.value
        });
    }
    console.log(info);
}
main();
