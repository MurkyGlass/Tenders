
    console.log('перешел па')
    document.getElementById('accountrevokebtn').addEventListener('click',Revoke)
async function Revoke() {
    const response = await fetch("/auth/revoke");
    alert("запрос отправлен")
    alert(response.status)
    if (!response.ok){
        alert("Ошибка выхода из аккаунта, попробуйте снова, через какое то время")   
        return      
    }
    window.location.href = "/main";
}