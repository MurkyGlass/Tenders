function applySort(sortValue) {
    const tendersList = document.getElementById('tendersList');
    tendersList.style.opacity = '0.5';
    
    const url = new URL(window.location);
    url.searchParams.set('sort', sortValue);
    
   window.location.href = url;
}


