<script lang="ts">
  import "smelte/src/tailwind.css";
  import { TextField, Checkbox, Chip } from "smelte";
  import Service, { Jq } from "./service";
  import { debounce } from "lodash-es";

  export let global: any;
  const jq = new Jq();
  let result = "";
  let jqVersion = "...";

  const service = new Service();
  service.getJqVersion().then((ver) => {
    jqVersion = ver;
  });
  service.getJqInput().then((json) => {
    jq.j = json;
  });

  $: {
    onChangeJq(jq);
  }

  function onClickDocumentation() {
    window.open("https://stedolan.github.io/jq/manual/", "_blank");
  }

  const onChangeJq = debounce((jq) => {
    service.runJq(jq).then((output) => {
      result = output;
    });
  }, 250);
</script>

<main>
  <div class="nav container max-w-none p-2">
    <div class="inline-block align-middle">
      <a href="/" class="navbar-brand"
        ><img src="images/logo.png" alt="jqplay" /></a
      >
    </div>
    <div class="inline-block text-sm align-middle ml-4">
      A playground for <a
        href="https://stedolan.github.io/jq/"
        class="navbar-link">jq</a
      >
      {jqVersion}
    </div>
    <p class="inline-block docs">
      <Chip icon="article" outlined on:click={onClickDocumentation}
        >documentation</Chip
      >
    </p>
  </div>

  <div class="main container max-w-none p-2">
    <div class="grid grid-cols-2 gap-4 h-600px">
      <div class="inputs">
        <h6>JQ Filter</h6>
        <TextField bind:value={jq.q} outlined />

        <h6>JSON</h6>
        <TextField bind:value={jq.j} textarea outlined />
      </div>

      <div class="outputs">
        <h6 class="">Result</h6>
        <div class="flex">
          <Checkbox name="compact" label="Compact Output" value="1" />
          <Checkbox label="Null Input" />
          <Checkbox label="Raw Input" />
          <Checkbox label="Raw Output" />
          <Checkbox label="Slurp" />
        </div>
        <TextField
          bind:value={result}
          label="Output"
          placeholder=""
          textarea
          outlined
          readonly
        />
      </div>
    </div>
  </div>
</main>

<style>
  main,
  div.main {
    @apply h-full;
    .grid {
      @apply h-full;
    }
  }

  :global {
    .outputs {
      textarea {
        height: 220px;
      }
    }
  }

  .nav {
    background-color: #222;
    color: #9d9d9d;
  }

  h6 {
    font-weight: bolder;
  }

  .main {
    height: 600px;
  }
</style>
