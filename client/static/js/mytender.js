window.addEventListener('pageshow', function(event) {
    if(event.persisted){
        this.window.location.href = "/protected/lk/tenders"
    }
});