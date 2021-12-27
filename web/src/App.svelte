<script lang="ts">
  import "smelte/src/tailwind.css";
  import { TextField, Checkbox, Chip, ProgressLinear, Button } from "smelte";
  import Service, { Jq } from "./service";
  import { debounce } from "lodash-es";

  import { json as langJson } from "@codemirror/lang-json";
  import CodeMirror from "./CodeMirror.svelte";
  import { onMount } from "svelte";
  import { samplesLeft, samplesRight } from "./samples";

  export let global: any;
  const jq = new Jq();
  let result = "";
  let jqVersion = "...";
  let commandLine = "";

  const langs = [langJson()];
  let jsonInputHeight = "500px";

  const slurp = jq.getOpt("slurp");
  const nullInput = jq.getOpt("null-input");
  const compactOutput = jq.getOpt("compact-output");
  const rawInput = jq.getOpt("raw-input");
  const rawOutput = jq.getOpt("raw-output");

  const service = new Service();
  service.getJqVersion().then((ver) => {
    jqVersion = ver;
  });

  onMount(() => {
    // load initial json, if avail
    service.getJqInput().then((json) => {
      if (!json) {
        json = "";
      }
      jq.j = json;
    });
  });

  function onClickDocumentation() {
    window.open("https://stedolan.github.io/jq/manual/", "_blank");
  }

  // run jq on input/filter/option changes
  $: {
    onChangeJq(jq);
  }
  const onChangeJq = debounce((jq) => {
    startProgressBar();
    service
      .runJq(jq)
      .then((output) => {
        progress = -1;
        clearTimeout(progressTimeoutId);
        result = output;
      })
      .then(() => {
        service.getJqCommandLine(jq).then((output) => {
          commandLine = output;
        });
      });
  }, 250);

  // progressbar
  let progress = -1;
  let progressTimeoutId = null;
  function startProgressBar() {
    $: {
      if (progress < 0 && progressTimeoutId !== null) {
        clearTimeout(progressTimeoutId);
      }
    }
    progress = 0;
    function next() {
      progressTimeoutId = setTimeout(() => {
        if (progress < 0) {
          return;
        }
        if (progress === 100) {
          progress = 0;
        }
        progress += 1;
        next();
      }, 25);
    }
    next();
  }

  // load cheatsheet sample
  function loadSample(sample) {
    jq.j = sample.input_j;
    jq.q = sample.input_q;
  }
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

  {#if progress >= 0}
    <ProgressLinear {progress} />
  {/if}

  <div class="main container max-w-none p-2">
    <div class="grid grid-cols-2 gap-4 h-600px">
      <div class="inputs">
        <h6>JQ Filter</h6>
        <TextField bind:value={jq.q} outlined />

        <h6>JSON</h6>
        <CodeMirror
          class="json_input"
          bind:value={jq.j}
          {langs}
          --cm-height={jsonInputHeight}
        />
      </div>

      <div class="outputs">
        <h6 class="">Result</h6>
        <div class="jq_options flex">
          <Checkbox
            label="Compact Output"
            bind:checked={compactOutput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Null Input"
            bind:checked={nullInput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Raw Input"
            bind:checked={rawInput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Raw Output"
            bind:checked={rawOutput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Slurp"
            bind:checked={slurp.enabled}
            on:change={onChangeJq(jq)}
          />
        </div>
        <CodeMirror
          class="json_output"
          bind:value={result}
          {langs}
          readOnly={true}
        />
      </div>
    </div>

    <div class="commandline mt-8 border">
      <h6 class="text-center pt-4 pb-4">Command Line</h6>
      <div class="text-center mt-2 mb-4"><code>{commandLine}</code></div>
    </div>

    <div class="cheatsheet border">
      <h6 class="text-center pt-4 pb-4">Cheatsheet</h6>
      <div class="grid grid-cols-2 gap-4 p-4">
        <table class="table-auto">
          <tbody>
            {#each samplesLeft as sample}
              <tr class="border-t">
                <td>
                  <code>{sample.code}</code>
                </td>
                <td>{sample.text}</td>
                <td>
                  <Button
                    icon="assignment"
                    text
                    light
                    flat
                    on:click={() => loadSample(sample)}
                  />
                </td>
              </tr>
            {/each}
          </tbody>
        </table>

        <table class="table-auto">
          <tbody>
            {#each samplesRight as sample}
              <tr class="border-t">
                <td>
                  <code>{sample.code}</code>
                </td>
                <td>{sample.text}</td>
                <td>
                  <Button icon="assignment" text light flat />
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</main>

<style>
  main,
  div.main {
    /* @apply h-full; */
    .grid {
      /* @apply h-full; */
    }
  }

  :global {
    .outputs {
      textarea {
        height: 220px;
      }
    }
  }

  :global {
    .commandline input {
      font-family: "SF Mono", "DejaVu Sans Mono", Menlo, Monaco, Consolas,
        Courier, monospace !important;
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

  /* reduce checkbox padding */
  :global {
    .jq_options {
      label {
        padding-left: 0 !important;
      }
    }
  }

  .cheatsheet,
  .commandline {
    h6 {
      background-color: #f5f5f5;
    }
  }

  code {
    color: #c7254e;
    background-color: #f9f2f4;
    border-radius: 4px;
  }

  :global {
    .cheatsheet {
      th:first-child,
      td:first-child {
        border-right-width: 0;
      }
    }
  }
</style>
