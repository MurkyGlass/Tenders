
async function Search(){
    searcher = document.getElementById('searchbox')
    sorter = document.getElementById('sortSelect')
    const url = new URL(window.location);
    url.searchParams.set('sort', sorter.value);
    url.searchParams.set('search',searcher.value)
    await UpdateCompanies(url)
        
}
async function UpdateCompanies(url) {
    const response = await fetch(url)
        if(response.ok){
            const fullHtml = await response.text();
            
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = fullHtml;
            
            const newElement = tempDiv.querySelector('#tendersList');
            
            if (newElement) {
                const currentElement = document.querySelector('#tendersList');
                
                if (currentElement) {
                    currentElement.replaceWith(newElement);
                }
            }
        }
}
async function applySort(sortValue) {
    searcher = document.getElementById('searchbox')
    const url = new URL(window.location);
    url.searchParams.set('sort', sortValue);
    url.searchParams.set('search',searcher.value)
  
    await UpdateCompanies(url)
}



