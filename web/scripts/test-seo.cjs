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
  const file = path.resolve(__dirname, '..', relative);
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

  let count = 0;
  mock.set('react', { ...require('react'), cache: fn => fn });
  mock.set('./ArticleList', { default: () => null });
  mock.set('@/src/api/category', { getCategoryNameStringByRoute: async () => 'Backend' });
  mock.set('@/src/api/tags', { getTagNameByRoute: async () => 'Go 1.23' });
  mock.set('@/src/api/posts', { getPostList: async () => ({ totalCount: count, totalPages: Math.ceil(count / 10), list: [] }) });
  mock.set('@/src/api/config', { getCommonConfig: async () => ({seo_meta:{title:'Blog'},website_meta:{website_name:'Blog'}}) });
  mock.set('@/src/api/stats', {});
  mock.set('@/src/utils/publicUrl', {});
  mock.set('@/src/utils/seo', { getSiteUrl: () => new URL('https://example.com') });
  const { archiveMetadata } = load('src/components/ArchivePage.tsx');
  let metadata = await archiveMetadata('tags','go1.23',undefined,{});
  assert.deepEqual(metadata.robots,{index:false,follow:true});
  assert.equal(metadata.alternates.canonical, '/tags/go1.23');
  count=21;
  metadata = await archiveMetadata('tags','go1.23','2',{});
  assert.equal(metadata.robots,undefined);
  assert.equal(metadata.alternates.canonical,'/tags/go1.23/page/2');
  await assert.rejects(archiveMetadata('tags','go1.23','4',{}), e => e.status === 404);
  const config=load('next.config.ts').default;
  const redirects=await config.redirects();
  assert(!redirects.some(r => r.source==='/about'));
  assert(redirects.some(r => r.source==='/posts/about-me' && r.destination==='/about-me' && r.permanent));
  console.log('SEO regression checks passed: pagination, redirects, canonical, empty archive recovery, and 404 versus service failures.');
})().catch(error => { console.error(error); process.exitCode=1; });
