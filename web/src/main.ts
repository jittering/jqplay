import App from './App.svelte';
import Global from './Global.svelte';

const app = new App({
  target: document.body,
  props: {
    global: Global,
  }
});

export default app;
