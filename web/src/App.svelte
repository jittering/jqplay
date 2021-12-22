<script lang="ts">
  import { TextField, Checkbox, Chip } from "smelte";
  import "smelte/src/tailwind.css";
  export let global: any;

  function onClickDocumentation() {
    window.open("https://stedolan.github.io/jq/manual/", "_blank");
  }
</script>

<main>
  <div class="nav container max-w-none p-2">
    <div class="inline-block">
      <a href="/" class="navbar-brand"
        ><img src="images/logo.png" alt="jqplay" /></a
      >
    </div>
    <div class="inline-block">
      A playground for <a
        href="https://stedolan.github.io/jq/"
        class="navbar-link">jq</a
      >
      TODO_JQ_VER <!-- TODO: -->
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
        <TextField outlined />

        <h6>JSON</h6>
        <TextField textarea outlined />
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
        <TextField label="Output" placeholder="" textarea outlined />
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
