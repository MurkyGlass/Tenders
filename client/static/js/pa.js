
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
        }
        EditformShow()
    });

function EditformShow(){
    document.getElementById('editModal').classList.add("show")
}
function EditformClose(){
    document.getElementById('editModal').classList.remove("show")
}
async function Revoke() {
    const response = await fetch("/auth/revoke");
    alert("запрос отправлен")
    alert(response.status)
    if (!response.ok){
        alert("Ошибка выхода из аккаунта, попробуйте снова, через какое то время")   
        return      
    }
    window.location.replace("/main");
}