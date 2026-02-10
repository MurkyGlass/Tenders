document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('accountrevoke-btn').addEventListener('click',Revoke) 
});

async function Revoke() {
    const response = await fetch("/auth/revoke");
    if (!response.ok){
        alert("Ошибка выхода из аккаунта, попробуйте снова, через какое то время")   
        return      
    }
    window.location.href = "/main";
}