window.addEventListener('pageshow', function(event) {
    if(event.persisted){
        window.location.href = this.window.location.pathname
    }
});
async function CompletTender(id) {
    const response = await fetch(window.location.pathname + "/offers/"+ id, {//
            method: 'POST',
        });
        if (response.ok){
            alert("Победитель успешно выбран")     
        }else{
            alert("Ошибка:"+ await response.text())
        }
        window.history.back()
}