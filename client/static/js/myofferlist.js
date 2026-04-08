window.addEventListener('pageshow', function(event) {
    if(event.persisted){
        this.window.location.href = "/protected/lk/offers"
    }
});
async function OpenMyOfferWindow(id){
    const resp = await fetch('/protected/lk/offers/'+id);
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                    OpenOfferWindow(id)
                }
                return
            } 
    window.location.href='/protected/lk/offers/'+id
}
async function Refresh(){
    await fetch("/auth/refresh");
}