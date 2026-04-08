window.addEventListener('pageshow', function(event) {
    if(event.persisted){
        this.window.location.href = "/protected/lk/tenders"
    }
});
async function viewMyTenderDetails(id){
    const resp = await fetch('/protected/lk/tenders/'+id);
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                    CreateOfferWindow(id)
                }
                if(resp.status == 409){
                    alert("Отказано в доступе")
                }
                return
            } 
    window.location.href='/protected/lk/tenders/'+id
}
async function Refresh(){
    await fetch("/auth/refresh");
}