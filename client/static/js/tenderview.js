async function CreateOfferWindow(id){
    const resp = await fetch('/protected/tenders/'+id+'/offer/create');
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                    CreateOfferWindow(id)
                }
                return
            } 
    window.location.href='/protected/tenders/'+id+'/offer/create'
}
async function OpenOfferWindow(id){
    const resp = await fetch('/protected/offers/'+id);
            if(!resp.ok){
                if(resp.status == 401){
                    document.getElementById('loginModal').classList.add("show")
                }else if(resp.status == 403){
                    await Refresh()
                    OpenOfferWindow(id)
                }
                return
            } 
    window.location.href='/protected/offers/'+id
}
async function Refresh(){
    await fetch("/auth/refresh");
}