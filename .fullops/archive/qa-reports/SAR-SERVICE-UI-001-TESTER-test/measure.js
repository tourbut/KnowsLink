(()=>{document.querySelectorAll('details').forEach(d=>d.open=true);
const iw=innerWidth;
const over=[...document.body.querySelectorAll('*')].filter(e=>{const b=e.getBoundingClientRect();return b.width>0&&(b.right>iw+1||b.left<-1)&&!e.closest('.skip')}).map(e=>e.tagName+'.'+e.className).slice(0,5);
const small=[...document.querySelectorAll('button,input:not([type=hidden]),select,summary,.toc a,a.button')].filter(e=>{const b=e.getBoundingClientRect();return b.width>0&&b.height<43.5}).map(e=>e.tagName+':'+(e.textContent||e.name).trim().slice(0,20)+':'+Math.round(e.getBoundingClientRect().height)).slice(0,5);
const unl=[...document.querySelectorAll('input:not([type=hidden]),select')].filter(e=>!e.id||!document.querySelector('label[for="'+e.id+'"]')).filter(e=>!e.closest('label')&&!e.getAttribute('aria-label')).map(e=>e.name).slice(0,5);
return JSON.stringify({iw,sw:document.documentElement.scrollWidth,hscroll:document.documentElement.scrollWidth>iw,over,small,unl,h1:document.querySelectorAll('h1').length,main:document.querySelectorAll('main').length,forms:document.forms.length,lang:document.documentElement.lang})})()
