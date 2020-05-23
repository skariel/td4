<script>
	import { onMount, onDestroy } from 'svelte';
	import { get, getUser } from './utils';

	let user      = null;
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
			.then((r)=>{solutions=r.data;})
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

    .displayname {

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

</style>

<title>test {test_id}</title>

<div class="top">
	<img class="avatar" src={test.avatar} alt="avatar"/>
	<h4 class="displayname">{test.display_name}</h4>
	<h4 class="testid">{test.id}</h4>
</div>

<pre class="code">
	<code>
		{test.code}
	</code>
</pre>

<!-- <h4 class="updated">updated: {test.ts_updated}</h4> -->

<a href={'/new_solution?test_id='+test_id}>add solution</a>


<h4>solutions:</h4>

{#each solutions as s }
	<div  style="margin:15px;">
		<h5>s:</h5>
		<pre>
			<code>
				{s.code}
			</code>
		</pre>
		<a href="/test/{test_id}/code/{s.id}">code id: {s.id}</a>
	</div>
{/each}
