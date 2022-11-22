<template>
  <div>
    <div>
      <h1>Hello {{ user.Username }}</h1>
      <div>
        <p>email</p>
        <input v-model="email" />
      </div>
      <button @click="update">Update</button>
      <button @click="test">test</button>
    </div>
    <div>
      <div class="exercices">
        <div v-for="(item, index) in statement" :key="index" class="child">
          <div v-html="item"></div>
          <div
            style="display: flex; justify-content: center; align-items: center"
          >
            <span>
              <h6>C</h6>
              <input @change="uploadFile($event, 'c', index)" type="file" />
            </span>
            <span>
              <h6>python</h6>
              <input @change="uploadFile($event, 'python', index)" type="file" />
            </span>
            <h6>Grads : {{ grads }} statement: {{ state }}</h6>
            <button @click="candyfresh('c')">Refresh</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "@vue/reactivity";
import axios from "axios";

var user = ref({});
var email = ref();
const grads = ref(0);
const state = ref("not submitted");
const exercices = ref();
const file = ref(null);
const statement = ref({
  fusion: `
      <h1 id="tri-fusion">Tri Fusion</h1>
      <pre><code>ALGORITHME TriFusion(T, gauche, droite);
        T: tableau de valeurs;
        gauche,droite: entier;
        centre: entier;
      DEBUT  
        SI (gauche &lt; droite) ALORS
            centre ← (gauche + droite) / 2;
            TriFusion (T, gauche, centre);
            TriFusion (T, centre + 1, droite);
            FUSIONNER(T, gauche, centre, droite);
        FSI
      FIN.
      </code></pre>
      <p><strong>Rendu utilisé la fonction suivant</strong></p>
      <pre><code class="lang-c">void afficher(int tab[], int n)
      {
          for (int i = 0; i &lt; n; i++)
              printf(&quot;%d\n&quot;, tab[i]);
      }
      </code></pre>
      <p>trier les tableaux {64, 25, 12, 22, 11} et {18, 88, 144, 45, 52, 31, 59, 108, 68, 124}</p>
      </p>
    `,
  hanoï: `
    <h1 id="tour-hanoï">Tour Hanoï</h1>
    <p>tour de « départ » à une tour d&#39;« arrivée » en passant par une tour « intermédiaire », et ceci en un minimum de coups, tout en respectant les règles suivantes :</p>
    <p>on ne peut déplacer plus d&#39;un disque à la fois ;
    on ne peut placer un disque que sur un autre disque plus grand que lui ou sur un emplacement vide.
    On suppose que cette dernière règle est également respectée dans la configuration de départ.</p>
    <p><strong> Rendu </srtong><br>
    Soit n le nième terme et X,Y ∈ {A,B,C} alors afficher </p>
    <pre><code class="language-txt">Disque n de X à Y
    </code></pre>
    `,
  expo: `
    <h1> Exponentiation Rapide</h1>
    <p>En mathématiques, plus précisément en arithmétique modulaire, l'exponentiation modulaire est un type d'élévation à la puissance (exponentiation) réalisée sur des entiers modulo un entier. Elle est particulièrement utilisée en informatique, spécialement dans le domaine de la cryptologie.
    </p><p>Etant donnés une base <i>b</i>, un exposant <i>e</i> et un entier non nul <i>m,</i> l'exponentiation modulaire consiste à calculer <i>c</i> tel que&#160;:
    </p>
    <dl><dd><span class="mwe-math-element"><span class="mwe-math-mathml-inline mwe-math-mathml-a11y" style="display: none;"><math xmlns="http://www.w3.org/1998/Math/MathML"  alttext="{\displaystyle c\equiv b^{e}{\pmod {m}}}">
      <semantics>
        <mrow class="MJX-TeXAtom-ORD">
          <mstyle displaystyle="true" scriptlevel="0">
            <mi>c</mi>
            <mo>&#x2261;<!-- ≡ --></mo>
            <msup>
              <mi>b</mi>
              <mrow class="MJX-TeXAtom-ORD">
                <mi>e</mi>
              </mrow>
            </msup>
            <mrow class="MJX-TeXAtom-ORD">
              <mspace width="1em" />
              <mo stretchy="false">(</mo>
              <mi>mod</mi>
              <mspace width="0.333em" />
              <mi>m</mi>
              <mo stretchy="false">)</mo>
            </mrow>
          </mstyle>
        </mrow>
        <annotation encoding="application/x-tex">{\displaystyle c\equiv b^{e}{\pmod {m}}}</annotation>
      </semantics>
    </math></span><img src="https://wikimedia.org/api/rest_v1/media/math/render/svg/4d786e34abdd71edd3a05ecf4430cf65c2174904" class="mwe-math-fallback-image-inline" aria-hidden="true" style="vertical-align: -0.838ex; width:17.826ex; height:2.843ex;" alt="c\equiv b^{e}{\pmod  {m}}"/></span></dd>
    <dd><span class="mwe-math-element"><span class="mwe-math-mathml-inline mwe-math-mathml-a11y" style="display: none;"><math xmlns="http://www.w3.org/1998/Math/MathML"  alttext="{\displaystyle 0\leq c&lt;m}">
      <semantics>
        <mrow class="MJX-TeXAtom-ORD">
          <mstyle displaystyle="true" scriptlevel="0">
            <mn>0</mn>
            <mo>&#x2264;<!-- ≤ --></mo>
            <mi>c</mi>
            <mo>&lt;</mo>
            <mi>m</mi>
          </mstyle>
        </mrow>
        <annotation encoding="application/x-tex">{\displaystyle 0\leq c&lt;m}</annotation>
      </semantics>
    </math></span><img src="https://wikimedia.org/api/rest_v1/media/math/render/svg/a9365f1c7a87728502c392dc8c7a40ebc01ebb48" class="mwe-math-fallback-image-inline" aria-hidden="true" style="vertical-align: -0.505ex; width:10.407ex; height:2.343ex;" alt="0\leq c&lt;m"/></span></dd></dl>
    <p>Par exemple, si <i>b</i> = 5, <i>e</i> = 3, et <i>m</i> = 13, le calcul de <i>c</i> donne 8.
    </p><p>Calculer l'exponentiation modulaire est considéré comme facile, même lorsque les nombres en jeu sont énormes. Au contraire, calculer le logarithme discret (trouver <i>e</i> à partir de <i>b</i>, <i>c</i> et <i>m</i>) est reconnu comme difficile. Ce comportement de fonction à sens unique fait de l'exponentiation modulaire une bonne candidate pour être utilisée dans les algorithmes de cryptologie.
    </p>
    <meta property="mw:PageProp/toc" />
   
`,
});



const props = defineProps({
  user: Object,
});

const uploadFile = (event, lang, index) => {
  console.log(index, lang);
  submitFile(event.target.files[0])
    .then((result) => {
      axios.post(
        "http://localhost:8080/exerices/correct",
        JSON.stringify({
          user_id: user.user_id,
          content: result,
          exo_id: index,
          exo_lang: lang
        }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      console.log();
    })
    .catch((e) => {});
};

const submitFile = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsText(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
  });
};

const test = async () => {
  const { data } = await axios.get(`http://localhost:8000/exercices`);
  console.log(data.Ok[0]);
}
const candyfresh = async () => {
  const user_id = props.user.Id
  ///const { data } = await axios.get(`http://localhost:8000/exercices/user/${user.value.user_id}`);
  //console.log(data);
};
const update = async () => {
  if (email.value) user.value.Email = email.value;

  axios.patch(
    "http://localhost:3000/user",
    {
      password: "",
      email: email.value,
    },
    {
      headers: { Authorization: "Bearer " + sessionStorage.token },
    }
  );
};
</script>

<style scoped lang="sass">
.exercices
  display: grid
  grid-template: repeat(1fr, 6)
  gap: 3rem
  margin-top: 2rem
.child
  grid-colomn: span 2
</style>