document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('accountBtn').addEventListener('click', LKbtn_Click);
    document.getElementById('loginSubmit').addEventListener('click', Authorization);
});
async function LKbtn_Click() {
            const resp = await fetch("../api/lk",{
                            method: 'GET',
                            headers:{
                                'Authorization':localStorage.getItem("access")
                            }
                        });
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                }
            }
       }
async function Authorization() {
    alert("нач")
    let login = document.getElementById("loginInput").value
    let password = document.getElementById("passwordInput").value
    if(login == "" || password == ""){
        alert("Заполните оба поля..")
        return
    }
    const loginForm = document.getElementById('loginForm');
    const formData = new FormData(loginForm);
  
    const response = await fetch('/auth/login', {
        method: 'POST',
        body: formData
    });
    if (response.ok){
        data = await response.json()
        let token = data.token_type +" "+ data.access_token
        localStorage.setItem("access",token)
        alert("Авторизация успешна")
        document.getElementById('loginModal').classList.remove("show")
    }
    if(response.status == 403){
        alert("Неправильный логин или пароль")        
    }
    await LKbtn_Click()
}
async function Refresh(){
    const response = await fetch("/auth/refresh",{
        method: 'GET',
        });
    if (response.ok){
        data = await response.json()
        let token = data.token_type + data.access_token
        localStorage.setItem("access",token)
                
    }else if(response.status == 401){
        localStorage.removeItem("access")
        await LKbtn_Click()
    }
    return
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