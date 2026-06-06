
    document.getElementById('accountrevokebtn').addEventListener('click',Revoke)
    document.getElementById('editbtn').addEventListener('click',EditformShow)
    document.getElementById('editClose').addEventListener('click',EditformClose)
    document.getElementById('editbackbtn').addEventListener('click',EditformClose)

window.addEventListener('pageshow', function(event) {
    if(event.persisted){
        this.window.location.href = "/protected/lk"
    }
});

document.getElementById('editForm').addEventListener('submit', async function(event) {

        event.preventDefault()
        const formData = new FormData(this);
  
        const response = await fetch('/protected/lk/edit', {
            method: 'POST',
            body: formData
        });
        if (response.ok){
            alert("Изменение успешно")
            window.location.href = "/protected/lk"
            EditformClose()
            return
        }else{
            alert("Ошибка:"+await response.text())
        }
        EditformShow()
    });
document.getElementById('roleCreateForm').addEventListener('submit', async function(event) {

        event.preventDefault()
        const formData = new FormData(this);
  
        const response = await fetch('/protected/lk/company/role/create', {
            method: 'POST',
            body: formData
        });
        if (response.ok){
            alert("Должность успешно создана!")
            window.location.href = "/protected/lk"
            RoleformClose()
            return
        }else{
            alert("Ошибка:"+await response.text())
        }
        RoleformShow()
    });
        async function RoleformShow(){
             const response = await fetch('/protected/lk/company/role/create', {
            method: 'POST'
            });
            if (response.status === 409){
                alert("Ошибка:"+await response.text())
                return
            }
            document.getElementById('roleCreateModal').classList.add("show")
        }
        function RoleformClose(){
            document.getElementById('roleCreateModal').classList.remove("show")
        }
async function EditformShow(){
    const response = await fetch('/protected/lk/edit', {
            method: 'POST'
            });
            if (response.status === 409){
                alert("Ошибка:"+await response.text())
                return
            }
    document.getElementById('editModal').classList.add("show")
}
function EditformClose(){
    document.getElementById('editModal').classList.remove("show")
}
async function Revoke() {
    const response = await fetch("/auth/revoke");
    alert("запрос отправлен")
    if (!response.ok){
        alert("Ошибка выхода из аккаунта, попробуйте снова, через какое то время")   
        return      
    }
    window.location.replace("/main");
}