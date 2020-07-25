<svelte:head>
	<link rel="stylesheet"
      href="//cdn.jsdelivr.net/gh/highlightjs/cdn-release@10.1.2/build/styles/default.min.css">
	<script src="//cdn.jsdelivr.net/gh/highlightjs/cdn-release@10.1.2/build/highlight.min.js"></script>

	<title>Test {test_id}</title>
</svelte:head>

<script>
	import { onMount, onDestroy, afterUpdate } from 'svelte';
	import { goto } from '@sapper/app';

	import { get, del, getUser, loginpath, start_invalidate_cache, timeSince } from './utils';
	import SolutionCard from '../components/SolutionCard.svelte'

	let user      = {};
    let page      = 0;
    let test_id   = 0;
	let test      = [];
	let solutions = [];
	let loading   = 2;
	let delete_count = 3;
	const delete_msg = [
		'permanently delete test, final button!',
		'press again to delete, just to be sure',
		'click it one more time!',
		'delete test',
	]
	let curr_delete_msg = delete_msg[3];

	afterUpdate(()=>{
		document.querySelectorAll('pre code').forEach((block) => {
			hljs.highlightBlock(block);
		})
	})

	onMount(async ()=>{
		delete_count = 3;
        user = getUser();
		window.addEventListener("locationchange", load_data);
		load_data()
	});

	onDestroy(async ()=>{
        if (process.browser == true) {
            window.removeEventListener("locationchange", load_data);
        }
	})

	function params(page) {
		return `/test?id=${test_id}&page=${page}`
	}

	function load_data() {
		loading = 2;
		const pname_wo_trailin_slash = window.location.pathname.replace(/\/+$/, '');
        if (pname_wo_trailin_slash != '/test') {
            return
        }
		const url = new URL(location)
		page = url.searchParams.get("page")
		if (page == null) {
			page = "0";
		}
		page = parseInt(page)
        test_id = url.searchParams.get("id")
		get(user, 'test/'+test_id)
			.then((r)=>{
				test=r.data;
				loading -= 1;
			  });
		get(user, 'solutions_by_test/'+test_id+'/'+page*10)
			.then((r)=>{solutions=r.data; loading -= 1;})
	}

	async function delete_test() {
		if (delete_count > 0) {
			delete_count -= 1;
			curr_delete_msg = delete_msg[delete_count];
			return
		}
		const res = await del(user, 'delete_test/'+test.id);
		if (res.status == 200) {
			start_invalidate_cache();
			goto('/')
		}
	}

</script>


<style>
	.solutions {
 		column-count: 2;
		column-gap: 1em;
	}

	@media (max-width:650px) {
		.solutions {
			column-count: 1;
			column-gap: 1em;
		}
	}


	.solutioncard {
		display: inline-block;
		margin-top: 10px;
		width: 100%;
	}

    .top {
        display:       flex;
    }

    .teststat {
        display:       flex;
		margin-top:    0px; 
    }

    .avatar {
        width:    40px;
        height:   40px;
        margin-right:20px;
    }

	.testid {
		margin-left: auto;
		font-size: 150%;
		margin-top: 7px;
	}

	.title {
		display: flex;
		align-items: center;
	}

	.title a {
		margin-left: auto;
	}

	.bottom {
		display: flex;
		justify-content: center;
		margin-top: 20px;
	}

    .filler {
        width: 20px;
    }

	.updated {
	    grid-column: 1 / 5;
	    grid-row: 4;
	    font-size: 80%;
		margin-top: -5px;
  	}



</style>

{#if loading>0}
	<h1>Loading...</h1>
{:else}

  <div class="top">
    <img class="avatar" src={test.avatar} alt="avatar" />
    <div style="display:flex; flex-direction: column; width: 100%">
      <h4 class="displayname">{test.display_name}</h4>
      <h4 class="updated">{timeSince(new Date(test.ts_updated))}</h4>
    </div>
    <h4 class="testid"><a href={'/test?id=' + test.id + '&page=0'}>#{test.id}</a></h4>
  </div>


	<h3>{test.title}</h3>
	<p style='color:#777777;'>{test.descr}</p>

	{#if user.display_name==test.display_name}
		<button on:click={goto("/test_edit?id="+test.id)}>Edit Test</button>
		<button on:click={delete_test}>{curr_delete_msg}</button>
	{/if}

	<pre class="code">
		<code class='lang-python'>
			{test.code}
		</code>
	</pre>

	<div class="title">
		{#if solutions.length > 0}
			<h1>All solutions</h1>
		{:else}
			{#if page == 0}
				<h1>No solutions yet!</h1>
			{:else}
				<h1>No solutions in this page!</h1>
			{/if}
		{/if}
		{#if user['avatar'] != null}
			<a href={"/new_solution?test_id="+test_id}>Add Solution</a>
		{:else}
			<a href={loginpath()}>Login to add a solution</a>
		{/if}
	</div>

	{#if solutions.length > 0}
		<div class="teststat">
			<h4>fail: {test.total_fail}</h4>
			<h4 style="margin-left:10px;">pass: {test.total_pass}</h4>
			<h4 style="margin-left:10px;">pending: {test.total_pending}</h4>
			<h4 style="margin-left:10px;">wip: {test.total_wip}</h4>
		</div>
	{/if}

	<div class="solutions">
		{#each solutions as s }
			<div class="solutioncard">
				<SolutionCard solution={s} />
			</div>
		{/each}
	</div>

	<div class="bottom">
		{#if page==0}
			{#if solutions.length == 10}
				<a href="{params(page+1)}">Next Page</a>
			{/if}
		{:else if solutions.length < 10}
			<a href="{params(page-1)}">Prev Page</a>
		{:else}
			<a href="{params(page-1)}">Prev Page</a>
			<div class="filler" />
			<a href="{params(page+1)}">Next Page</a>
		{/if}
	</div>


{/if}