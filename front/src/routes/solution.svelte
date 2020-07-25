<svelte:head>
	<link rel="stylesheet"
		href="//cdn.jsdelivr.net/gh/highlightjs/cdn-release@10.1.2/build/styles/default.min.css">
	<script src="//cdn.jsdelivr.net/gh/highlightjs/cdn-release@10.1.2/build/highlight.min.js"></script>

	<title>Solution {solution_id}</title>
</svelte:head>

<script>
	import { onMount, afterUpdate } from 'svelte';
	import { get, del, getUser, start_invalidate_cache, timeSince } from './utils';
	import { goto } from '@sapper/app';

    let user          = {};
    let solution_id   = 0;
	let solution      = [];
	let results       = [];
	let loading       = 2;
	let delete_count = 3;
	const delete_msg = [
		'permanently delete test, final button!',
		'press again to delete, just to be sure',
		'click it one more time!',
		'delete test',
	]
	let curr_delete_msg = delete_msg[3];
    let back_color = '';



	afterUpdate(()=>{
		document.querySelectorAll('pre code').forEach((block) => {
			hljs.highlightBlock(block);
		})
	})

	onMount(async ()=>{
        user = getUser();
		load_data()
	});

	async function load_data() {
		loading = 2;
		const url = new URL(location)
        solution_id = url.searchParams.get("id")
		let r = await get(user, 'solution/'+solution_id)
		solution=r.data;
		r = await get(user, 'results_by_run/'+solution.run_id)
		results=r.data;
		loading = 0;
		setBackColor();
	}

	async function delete_solution() {
		if (delete_count > 0) {
			delete_count -= 1;
			curr_delete_msg = delete_msg[delete_count];
			return
		}
		const res = await del(user, 'delete_solution/'+solution.id);
		if (res.status == 200) {
			start_invalidate_cache();
			goto("/test?id="+solution.test_code_id);
		}
	}

	function setBackColor() {
		if (solution.status=='fail') {
			back_color = 'failed_back';
		}
		else if (solution.status=='pass') {
			back_color = 'passed_back';
		}
	}

</script>

<style>
    .top {
        display:       flex;
    }
    .avatar {
        width:    40px;
        height:   40px;
        margin-right:20px;
    }

    .solutionid {
		margin-left: auto;
		font-size: 100%;
		margin-top: 0px;
    }

	.title {
		display: flex;
		align-items: center;
	}
	.updated {
		font-size: 80%;
		margin-top: -5px;
	}
    .failed_back {
        background-color: rgb(130,55,55);
    }

    .passed_back {
        background-color: rgb(55,130,55);
    }
	.status_light {
        margin-left:auto;
        border-radius: 17px;
        margin-bottom: 7px;
        width: 21px;
		height: 21px;
    }


</style>


{#if loading>0}
	<h1>loading...</h1>
{:else}

  <div class="top">
    <img class="avatar" src={solution.avatar} alt="avatar" />
    <div style="display:flex; flex-direction: column; width: 100%">
      <h4 class="displayname">{solution.display_name}</h4>
      <h4 class="updated">{timeSince(new Date(solution.ts_updated))}</h4>
    </div>
    <div style="display:flex; flex-direction: column; width: 100%">
	    <h4 class="solutionid">solution #{solution.id}</h4>
	    <h4 class="solutionid"><a href={"/test?id="+solution.test_code_id}>test #{solution.test_code_id}</a></h4>
	</div>
  </div>

	<div style="display:flex; ">
		<h4>status: {solution.status} {#if solution.status=='stop'} (timeout) {/if}</h4>
		<div class="status_light {back_color}" />
	</div>


	{#if user.display_name==solution.display_name}
		<button on:click={goto("/solution_edit?id="+solution.id)}>Edit Code</button>
		<button on:click={delete_solution}>{curr_delete_msg}</button>
	{/if}

	<pre>
		<code>
			{solution.code}
		</code>
	</pre>

	{#if results.length > 0}
		{#if results[0].output.String.length > 0}
			<h2>output:</h2>
			<pre class="code">
				<code>
					{results[0].output.String}
				</code>
			</pre>
		{:else}
			<h2>No output</h2>
		{/if}
	{:else}
		<div class="title">
			<h2>No result yet</h2>
		</div>
	{/if}

{/if}