<script lang="ts">
  import { Button } from "smelte";

  export let label: string;
  export let collapsible = false;
  export let expanded = true;
</script>

<div class="flex-initial p-4">
  <div class="panel border">
    <div class="header text-center pt-4 pb-4">
      {#if collapsible}
        <Button
          classes="float-right -mt-4"
          icon={expanded ? "unfold_less" : "unfold_more"}
          text
          light
          flat
          on:click={() => {
            expanded = !expanded;
          }}
        />
      {/if}
      <h6>{label}</h6>
    </div>

    <div class="mt-2 mb-4">
      {#if !collapsible || expanded}
        <slot />
      {/if}
    </div>
  </div>
</div>

<style>
  .panel .header {
    background-color: #f5f5f5;
  }
</style>
