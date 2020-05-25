<script>
	import { onMount, onDestroy } from 'svelte';
	import { get, getUser, loginpath } from './utils';
	import SolutionCard from '../components/SolutionCard.svelte'

	let user      = {};
    let page      = 0;
    let test_id   = 0;
	let test      = [];
	let solutions = [];

	onMount(async ()=>{
        user = getUser();
		window.addEventListener("locationchange", load_data);
		load_data()
	});

	onDestroy(async ()=>{
        // TODO: otherwise running on node?! a bug?!
        if (process.browser == true) {
            window.removeEventListener("locationchange", load_data);
        }
	})

	function load_data() {
        if (window.location.pathname != '/test') {
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
			.then((r)=>{test=r.data;})
		get(user, 'solutions_by_test/'+test_id+'/'+page*10)
			.then((r)=>{solutions=r.data; console.log(r.data)})
	}
</script>

<style>
	.solutions {
 		column-count: 2;
		column-gap: 1em;
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
    }

	.code {
		background-color: #333333;
		min-height: 100px;
	}

	code {
		background-color: #00000000;
		color: antiquewhite;
	}
	.title {
		display: flex;
		align-items: center;
	}
	.title a {
		margin-left: auto;
	}

</style>

<svelte:head>
	<title>Test {test_id}</title>
</svelte:head>

<div class="top">
	<img class="avatar" src={test.avatar} alt="avatar"/>
	<h4 class="displayname">{test.display_name}</h4>
	<h4 class="testid">test {test.id}</h4>
</div>


<pre class="code">
	<code>
		{test.code}
	</code>
</pre>

<div class="title">
	{#if solutions.length > 0}
		<h1>All Solutions</h1>
	{:else}
		<h1>No Solutions Yet!</h1>
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

<!-- TODO: paging for solutions! -->
