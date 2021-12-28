<script lang="ts">
  import { Button } from "smelte";
  import { createEventDispatcher } from "svelte";
  import { samplesLeft, samplesRight } from "./samples";

  const dispatch = createEventDispatcher();
</script>

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
              on:click={() => dispatch("sample", { sample })}
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
