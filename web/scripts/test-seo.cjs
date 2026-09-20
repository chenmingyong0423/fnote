/* eslint-disable @typescript-eslint/no-require-imports -- Run TypeScript regression checks without a separate build. */
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const ts = require('typescript');

const modules = new Map();
const navigation = {
  notFound() { throw Object.assign(new Error('not found'), { status: 404 }); },
};
const mock = new Map([['next/navigation', navigation]]);
function load(relative) {
  let file = path.resolve(__dirname, '..', relative);
  if (!fs.existsSync(file) && file.endsWith('.ts')) file += 'x';
  if (modules.has(file)) return modules.get(file).exports;
  const mod = { exports: {} }; modules.set(file, mod);
  const code = ts.transpileModule(fs.readFileSync(file, 'utf8'), {
    compilerOptions: { target: ts.ScriptTarget.ES2017, module: ts.ModuleKind.CommonJS, jsx: ts.JsxEmit.ReactJSX, esModuleInterop: true },
  }).outputText;
  const localRequire = (name) => {
    if (mock.has(name)) return mock.get(name);
    if (name.startsWith('@/')) return load(name.slice(2) + (name.endsWith('Page') ? '.tsx' : '.ts'));
    if (name.startsWith('.')) return load(path.relative(path.resolve(__dirname, '..'), path.resolve(path.dirname(file), name + '.ts')));
    return require(name);
  };
  vm.runInThisContext('(function(require,module,exports){' + code + '\n})', { filename: file })(localRequire, mod, mod.exports);
  return mod.exports;
}

(async () => {
  const { resolvePagination, ensurePageExists } = load('src/utils/pagination.ts');
  const resolve = (page, query = {}) => resolvePagination('/tags/go1.23', page, query);
  assert.equal(resolve('2').pathname, '/tags/go1.23/page/2');
  assert.equal(resolve(undefined).pathname, '/tags/go1.23');
  for (const page of ['0', '1', '02', '001', '-1', 'abc', '1.5', '1e2', '9007199254740992', '']) {
    assert.throws(() => resolve(page), e => e.status === 404, page);
  }
  for (const pageSize of ['0', 'abc', '11', '1000000', ['10','20']]) {
    assert.throws(() => resolve(undefined, { pageSize }), e => e.status === 404);
  }
  for (const size of [5,10,20,50]) assert.equal(resolve('2', {pageSize:String(size)}).pageSize, size);
  assert.throws(() => resolve(undefined, {page:'2', keyword:'Go 语言', filter:'likes',pageSize:'20'}),
    e => e.status === 404);
  for (const page of ['1', '', ['1','2']]) {
    assert.throws(() => resolve(undefined, {page}), e => e.status === 404);
    assert.throws(() => resolve('2', {page}), e => e.status === 404);
  }
  assert.deepEqual(resolve('2', {filter:'likes',pageSize:'20'}),
    {page:2,pageSize:20,field:'likes',pathname:'/tags/go1.23/page/2'});
  assert.throws(() => ensurePageExists(9999, 3), e => e.status === 404);
  ensurePageExists(1, 0);

  const http = load('src/utils/http.ts');
  const originalFetch = global.fetch;
  try {
    const category = load('src/api/category.ts');
    const tags = load('src/api/tags.ts');
    global.fetch = async () => new Response(JSON.stringify({message:'missing'}), {status:404});
    assert.equal(await category.getCategoryNameStringByRoute('missing'), '');
    assert.equal(await tags.getTagNameByRoute('missing'), '');
    global.fetch = async () => new Response(JSON.stringify({message:'db unavailable'}), {status:500});
    await assert.rejects(tags.getTagNameByRoute('failed'), e => e instanceof http.HttpError && e.status === 500);
    await assert.rejects(category.getCategoryNameStringByRoute('failed'), e => e.status === 500);
    global.fetch = async () => { throw Error('offline'); };
    await assert.rejects(tags.getTagNameByRoute('failed'), e => e instanceof http.BackendUnavailableError);
  } finally { global.fetch = originalFetch; }

  const posts = load('src/api/posts.ts');
  try {
    global.fetch = async () => new Response(JSON.stringify({message:'missing'}), {status:404});
    assert.equal(await posts.getPostDetailOrNull('missing'), null);
    global.fetch = async () => new Response(JSON.stringify({message:'The postId does not exist.'}), {status:400});
    assert.equal(await posts.getPostDetailOrNull('legacy-missing'), null);
    global.fetch = async () => new Response(JSON.stringify({message:'Post not found.'}), {status:500});
    await assert.rejects(posts.getPostDetailOrNull('failed'), e => e.status === 500);
    global.fetch = async () => new Response(JSON.stringify({code:1,message:'database unavailable'}));
    await assert.rejects(posts.getPostDetailOrNull('failed'), /database unavailable/);
  } finally { global.fetch = originalFetch; }

  const seo = load('src/utils/seo.ts');
  const common = load('src/api/config.ts').DEFAULT_COMMON_CONFIG;
  const previousHost = process.env.BASE_HOST;
  process.env.BASE_HOST = 'https://example.com/';
  try {
    const metadata = seo.buildPageMetadata(common, {title:'内页标题',description:'内页摘要',pathname:'/navigation'});
    assert.equal(metadata.openGraph.url, 'https://example.com/navigation');
    assert.equal(metadata.twitter.title, metadata.title);
    assert.equal(metadata.twitter.description, '内页摘要');
    assert.equal(metadata.twitter.card, 'summary');
    const withImage = seo.buildPageMetadata(common, {title:'文章',description:'摘要',pathname:'/posts/test',image:'https://cdn.example.com/cover.jpg'});
    assert.equal(withImage.twitter.card, 'summary_large_image');
    assert.equal(withImage.openGraph.images[0].url, withImage.twitter.images[0].url);
    assert.equal(withImage.openGraph.images[0].alt, '文章');
    assert.equal(seo.getPostPath('中文 #?'), '/posts/%E4%B8%AD%E6%96%87%20%23%3F');
    assert.equal(seo.getPostPath('about-me'), '/about-me');
    const items = [{name:'首页',pathname:'/'},{name:'Go 语言',pathname:'/tags/go'}];
    const breadcrumb = seo.buildBreadcrumbJsonLd(items);
    assert.deepEqual(breadcrumb.itemListElement.map(item => item.position), [1, 2]);
    assert.equal(breadcrumb.itemListElement[1].item, 'https://example.com/tags/go');
    assert(!seo.serializeJsonLd({name:'</script><script>alert(1)</script>'}).includes('<'));
    const collection = seo.buildCollectionJsonLd('Go', '/tags/go', [{title:'教程',sug:'go 入门'}]);
    assert.equal(collection.mainEntity.itemListElement[0].url, 'https://example.com/posts/go%20%E5%85%A5%E9%97%A8');
    const Breadcrumbs = load('src/components/Breadcrumbs.tsx').default;
    const html = require('react-dom/server').renderToStaticMarkup(Breadcrumbs({items}));
    assert(html.includes('href="/"'));
    assert(html.includes('aria-current="page">Go 语言'));
    assert(html.includes('application/ld+json'));
  } finally {
    if (previousHost === undefined) delete process.env.BASE_HOST;
    else process.env.BASE_HOST = previousHost;
  }

  let count = 0;
  mock.set('react', { ...require('react'), cache: fn => fn });
  mock.set('./ArticleList', { default: () => null });
  mock.set('@/src/api/category', { getCategoryNameStringByRoute: async () => 'Backend' });
  mock.set('@/src/api/tags', { getTagNameByRoute: async () => 'Go 1.23' });
  mock.set('@/src/api/posts', { getPostList: async () => ({ totalCount: count, totalPages: Math.ceil(count / 10), list: [] }) });
  mock.set('@/src/api/config', { getCommonConfig: async () => ({seo_meta:{title:'Blog'},website_meta:{website_name:'Blog'}}) });
  mock.set('@/src/api/stats', {});
  const { archiveMetadata } = load('src/components/ArchivePage.tsx');
  let metadata = await archiveMetadata('tags','go1.23',undefined,{});
  assert.deepEqual(metadata.robots,{index:false,follow:true});
  assert.equal(metadata.alternates.canonical, '/tags/go1.23');
  count=21;
  metadata = await archiveMetadata('tags','go1.23','2',{});
  assert.equal(metadata.robots,undefined);
  assert.equal(metadata.alternates.canonical,'/tags/go1.23/page/2');
  assert.equal(metadata.twitter.title, metadata.title);
  for (const query of [{filter:'likes'}, {filter:'oldest'}, {pageSize:'20'}]) {
    const filtered = await archiveMetadata('tags','go1.23',undefined,query);
    assert.deepEqual(filtered.robots, {index:false,follow:true});
  }
  assert.equal((await archiveMetadata('tags','go1.23',undefined,{utm_source:'test'})).robots, undefined);
  await assert.rejects(archiveMetadata('tags','go1.23','4',{}), e => e.status === 404);
  const config=load('next.config.ts').default;
  const redirects=await config.redirects();
  assert(!redirects.some(r => r.source==='/about'));
  assert(redirects.some(r => r.source==='/posts/about-me' && r.destination==='/about-me' && r.permanent));
  console.log('SEO regression checks passed: pagination, redirects, page/social metadata, schema, server-rendered breadcrumbs, archive indexing, and 404 versus service failures.');
})().catch(error => { console.error(error); process.exitCode=1; });
