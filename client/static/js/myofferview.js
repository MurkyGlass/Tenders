window.addEventListener('pageshow', function(event) {
    if(event.persisted){
        window.location.href = this.window.location.pathname
    }
});
async function EditOfferWindow(id){
    const resp = await fetch('/protected/lk/offers/'+id+'/edit');
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                    EditOfferWindow(id)
                }
                if(resp.status == 409){
                    alert("Отказано в доступе")
                }
                return
            } 
    window.location.href='/protected/lk/offers/'+id+'/edit'
}
async function Refresh(){
    await fetch("/auth/refresh");
}