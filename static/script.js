const id = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);

// let source = new EventSource(`http://localhost:3500/events?id=${id}`);
// let source = new EventSource(`http://localhost:3500/events`);
let source = new EventSourcePolyfill(`http://localhost:3500/events`,
    {
        headers: {
            'Authorization': `Bearer ${id}`
        }
    });

source.addEventListener('open', (e) => {
    console.log("OPEN:", id);
});

source.addEventListener('test1', (e) => {
    console.log("SALUDAR:", e.data);
});

source.addEventListener('test2', (e) => {
    console.log("SALTAR:", e.data);
});
