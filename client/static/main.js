document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('accountBtn').addEventListener('click', LKbtn_Click);
    document.getElementById('loginForm').addEventListener('submit', async function(event) {

        event.preventDefault()
        const loginForm = document.getElementById('loginForm');
        const formData = new FormData(loginForm);
  
        const response = await fetch('/auth/login', {
            method: 'POST',
            body: formData
        });
        if (response.ok){
            alert("Авторизация успешна")
            document.getElementById('loginModal').classList.remove("show")
        }
        if(response.status == 403){
            alert("Неправильный логин или пароль")        
        }
        await LKbtn_Click()
    });
    document.getElementById('loginClose').addEventListener('click',LogFormClose);
    document.getElementById('passwordToggle').addEventListener('click',PassToggle);
    document.getElementById('regisBtn').addEventListener('click',RegistrShow);
    document.getElementById('registerClose').addEventListener('click',RegistrClose);
    document.getElementById('userPasswordToggle').addEventListener('click',UserPassToggle)
    document.getElementById('confirmPasswordToggle').addEventListener('click',ConfirmPassToggle)
    document.getElementById('registerForm').addEventListener('submit',async function(event){
        event.preventDefault()
        const Form = document.getElementById('registerForm');
        const formData = new FormData(Form);
  
        const response = await fetch('/main/registration', {
            method: 'POST',
            body: formData
        });
        if (response.ok){
            alert("Регистрация успешна")
            RegistrClose()
            return
        }

        RegistrShow()
    })
});
function ConfirmPassToggle() {
    let password = document.getElementById("confirmPassword")
    if(password.type == "password"){
        password.type = "text"
    }else if (password.type == "text"){
        password.type = "password"
    }
}
function UserPassToggle() {
    let password = document.getElementById("userPassword")
    if(password.type == "password"){
        password.type = "text"
    }else if (password.type == "text"){
        password.type = "password"
    }
}
function RegistrClose() {
    document.getElementById('registerModal').classList.remove("show")
}
function RegistrShow() {
    LogFormClose()
    document.getElementById('registerModal').classList.add("show")
}
function PassToggle() {
    let password = document.getElementById("passwordInput")
    if(password.type == "password"){
        password.type = "text"
    }else if (password.type == "text"){
        password.type = "password"
    }
}
function LogFormClose() {
    document.getElementById('loginModal').classList.remove("show")
}
async function LKbtn_Click() {
            const resp = await fetch("/protected/lk");
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                }
                return
            } 
            
        
            window.location.href = "/protected/lk";
}

async function Refresh(){
    await fetch("/auth/refresh");
    await LKbtn_Click()
}
document.querySelectorAll('.nav-menu a').forEach(link => {
    link.addEventListener('click', function(e) {
        if(!this.classList.contains('active')) {
            document.querySelectorAll('.nav-menu a').forEach(item => {
                item.classList.remove('active');
            });
            this.classList.add('active');
        }
    });
});