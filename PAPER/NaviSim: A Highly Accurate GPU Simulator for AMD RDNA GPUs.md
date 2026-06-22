NaviSim: A Highly Accurate GPU Simulator for AMD RDNA
GPUs
YuhuiBao YifanSun ZlatanFeric
bao.yu@northeastern.edu yifan@cs.wm.edu feric.z@northeastern.edu
NortheasternUniversity William&Mary NortheasternUniversity
Boston,MA,USA Williamsburg,VA,USA Boston,MA,USA
MichaelTianShen MicahWeston JoséL.Abellán
shen.mich@northeastern.edu weston.m@northeastern.edu jlabellan@ucam.edu
NortheasternUniversity NortheasternUniversity UniversidadCatólicadeMurcia
Boston,MA,USA Boston,MA,USA Murcia,Spain
TrinayanBaruah JohnKim AjayJoshi
tbaruah@amd.com jjk12@kaist.edu joshi@bu.edu
AMD KAIST BostonUniversity
SantaClara,CA,USA Daejeon,SouthKorea Boston,MA,USA
DavidKaeli
kaeli@ece.neu.edu
NortheasternUniversity
Boston,MA,USA
ABSTRACT TodemonstratethefullutilityoftheNaviSimsimulator,wecarry
AsGPUscontinuetogrowinpopularityforacceleratingdemanding outaperformancestudyoftheimpactofindividualRDNAfeatures,
applications,suchashigh-performancecomputingandmachine attemptingtounderstandbetterthedesigndecisionsbehindthese
learning,GPUarchitectsneedtodelivermorepowerfuldevices features.Wecarryoutanumberofexperimentstoisolateeach
withupdatedinstructionsetarchitectures(ISAs)andnewmicroar- RDNAfeatureandevaluateitsimpactonoverallperformance,as
chitecturalfeatures.TheintroductionoftheAMDRDNAarchitec- wellasdemonstratetheusabilityandflexibilityofNaviSim.
tureisoneexamplewheretheGPUarchitecturewasdramatically
CCSCONCEPTS
changed,modifyingtheunderlyingprogrammingmodel,thecore
architecture,andthecachehierarchy.Todate,nopublicly-available •Computersystemsorganization→Singleinstruction,mul-
simulatorinfrastructurecanmodeltheAMDRDNAGPU,prevent- tipledata;•Computingmethodologies→Modelingmethod-
ing researchers from exploring new GPU designs based on the ologies;Modelverificationandvalidation.
state-of-the-artRDNAarchitecture.
Inthispaper,wepresenttheNaviSimsimulator,thefirstcycle- KEYWORDS
levelGPUsimulatorframeworkthatmodelsAMDRDNAGPUs.
GPU,Simulation,ComputerArchitecture
NaviSimfaithfullyemulatesthenewRDNAISA.Weextensively
tune and validate NaviSim using several microbenchmarks and ACMReferenceFormat:
10fullworkloads.OurevaluationshowsthatNaviSimcanaccu- YuhuiBao,YifanSun,ZlatanFeric,MichaelTianShen,MicahWeston,
JoséL.Abellán,TrinayanBaruah,JohnKim,AjayJoshi,andDavidKaeli.
ratelymodeltheGPU’skernelexecutiontime,achievingsimilar
2022.NaviSim:AHighlyAccurateGPUSimulatorforAMDRDNAGPUs.In
performancetohardwareexecutionwithin9.92%(onaverage),as
InternationalConferenceonParallelArchitecturesandCompilationTechniques
measuredonanAMDRX5500XTGPUandanAMDRadeonPro
(PACT’22),October8–12,2022,Chicago,IL,USA.ACM,NewYork,NY,USA,
W6800GPU.
13pages.https://doi.org/10.1145/3559009.3569666
1 INTRODUCTION
Permissiontomakedigitalorhardcopiesofallorpartofthisworkforpersonalor
GPUshavebeenusedtoaccelerateawiderangeofmoderndata-
classroomuseisgrantedwithoutfeeprovidedthatcopiesarenotmadeordistributed
forprofitorcommercialadvantageandthatcopiesbearthisnoticeandthefullcitation centricapplications(e.g.,artificialintelligence[26],big-dataana-
onthefirstpage.CopyrightsforcomponentsofthisworkownedbyothersthanACM lytics[1],andhigh-performancecomputingworkloads[18]),lever-
mustbehonored.Abstractingwithcreditispermitted.Tocopyotherwise,orrepublish,
agingtheirever-increasingcomputingcapabilities[32].Withthe
topostonserversortoredistributetolists,requirespriorspecificpermissionand/ora
fee.Requestpermissionsfrompermissions@acm.org. continuousdemandforhigherperformance,GPUvendors(e.g.,
PACT’22,October8–12,2022,Chicago,IL,USA NVIDIAandAMD)havebeenpushingtheenvelopeofGPUperfor-
©2022AssociationforComputingMachinery.
manceineverynewGPUgeneration.WhilesomeGPUgenerations
ACMISBN978-1-4503-9868-8/22/10...$15.00
https://doi.org/10.1145/3559009.3569666 aremainly“specbumps”,otherGPUgenerationsintroducemajor
333

PACT’22,October8–12,2022,Chicago,IL,USA Bao,Sun,Feric,Shen,Weston,Abellán,Baruah,Kim,Joshi,andKaeli
DRAM
R9Nano’s L1$ L2$
RX 5500 XT’s L0$ L1$ L2$ DRAM
Figure 1: Comparing cache latencies between the GCN
and RDNA GPUs, while running the pointer chasing mi-
crobenchmark.Theresultssuggestthat,otherthanthepub-
licly announced changes in the specifications (e.g., cache
sizes),AMDhasmademanyunannounceddesignimprove-
ments(e.g.,reducingcachelatencies).
600
400
200
0
64 256 512 768 1024
Problem Size
)su(emiT
Todate,GPUarchitecturesimulatordevelopmentinacademia
hassignificantlylaggedbehindindustryschedules.Thestate-of-
the-artAMDGPUsimulatorsarestillmodelingtheAMDGCN3
architecture[15,33](releasedin2015),andnopubliclyavailable
simulatorsmodelthemorerecentAMDRDNAarchitectures.The
lackofup-to-datesimulatorinfrastructurescanimpactaGPUarchi-
tect’sabilitytoexplorenewinnovationstargetingnext-generation
GPUs.
ThislackofGPUsimulationtoolswithinacademiaispotentially
harmfultothecomputerarchitectureresearchcommunityinthe
longrun.Researchpapersoftendescribetheirresearchmethodol-
ogyusingsentencessuchas“wemodifysimulatorSthatisorigi-
nallyvalidatedagainstproductAtomodelproductBbychanging
thenumberofcoresandthecachesizes.”Typically,thisproduct
B is a few generations later than product A, and the study im-
plicitlyassumedthatthesimulatorcouldstillcorrectlymodelthe
newarchitecture.However,thenuancesinthearchitecturedesign
Simulation Native Execution
andparameterselectionmaysignificantlychangetheperformance
Turning point in simulation characteristics.Forexample,simplymodifyingMGPUSim[33](val-
but not in native execution idatedagainsta2015-releasedR9NanoGPU)tomodeltheRX5500
XTGPU(RDNAarchitecture),asshowninFigure2,canresultin
simulationerrorsaslargeas318%whenrunningtheBICGbench-
markfromPolyBench[30].Themajorinaccuraciesreportedabove
suggeststhatsimplyadjustingpubliclyknownparametersisinsuf-
Figure2:Acomparisonofasimulator’sreportedkernelexe- ficienttomodelabrandnewarchitectureandcanleadtowrong
cutiontimeandthenativeRX5500XTGPUexecutiontime conclusions.Acarefulredesignandre-calibrationofthesimula-
when running the BICG benchmark from PolyBench [30]. torinfrastructurearenecessarytoprovideatrustworthybaseline
The experiment is performed using MGPUSim [33], af- modelforthenextgenerationofGPUarchitectureresearch.
teronlymodifyingthepubliclyavailableparameters(core Toaddresstheseissues,wepresentNaviSim,aGPUsimulator
count,frequency,cachesizes,DRAMbandwidth,etc.)based thatmodelsAMDRDNAGPUs.NaviSimfaithfullyemulatesthe
on the validated R9 Nano model. Not only are the perfor- newAMDRDNAISAandproducesexactapplicationoutputsas
mancediscrepancieshigh(upto318%),buttheperformance recordedontheGPUhardware.UtilizingtheAkitasimulationen-
doesnotfollowthetrendsseenontheRX5500XThardware. gine[33],thesamesimulationenginethatdroveMGPUSim[33],
weenablemodularityandhigh-performanceparallelsimulation.
WeintegrateDRAMSim3[27]toaccuratelymodeldifferentDRAM
architecturaloverhauls,changingmanyaspectsoftheGPU’sorga- technologies,includingHBM[25]andGDDR[19].NaviSimisalso
nization.Forexample,in2019,AMDmovedawayfromtheGCN fullycompatiblewiththeDaisenGPUvisualizationframework[36],
architecture(usedbyalmostallAMDGPUssince2011)andstarted allowinguserstoeasilyunderstandtheGPU’sperformanceissues
torolloutanewRDNAarchitecturetobeusedontheir7nmGPUs. identifiedbythesimulator.WeextensivelyvalidateNaviSimagainst
AMD’sRDNAarchitecturewasamajorredesign,asAMDmodi- anAMDRX5500XTGPU(RDNAarchitecture).Theaveragesim-
fiednearlyeveryaspectoftheGPU’sarchitecture.Thisincluded ulationerroronexecutiontime,ascomparedtorealhardware,is
changestotheinstructionsetarchitecture(ISA),corearchitecture, lessthan10%acrossasuiteof10workloads.
andmemoryhierarchy.Inthecores,theRDNAarchitecturereduced Giventhatwehaveavalidatedsimulator,wethenuseNaviSim
thenumberofco-scheduledthreads(i.e.,thewavefront)from64 toanalyzetheeffectofdifferentRDNAfeaturesonaGPU’soverall
to32tocopewithahigherlevelofthreaddivergenceinmodern performance,attemptingtomakesenseofthedesigndecisions
workloads.Thememoryhierarchywasalsoextended,addingan behindthedesignoftheRDNAarchitecture.Wecarryoutanumber
extralayerofcachebetweentheoriginalL1andL2caches.This ofexperimentstoisolateeachRDNAfeatureandevaluateitsimpact.
reducedtheburdenoftheL2cachesandsimplifiedthemassive Theseexperimentsalsoserveascasestudies,demonstratingthe
L1toL2network.Inadditiontothesepubliclyknownchanges, usabilityandflexibilityofNaviSim.
enhancementswerealsomadeunderthehood.Forexample,aswe Insummary,thispapermakesthefollowingcontributions:
evaluatetheAMDRX5500XTGPU(RDNAarchitecture)against • WepresentNaviSim,ahighlyconfigurableandaccurateGPU
the R9 Nano GPU (GCN3 Architecture) with a pointer chasing simulatorthatmodelstheAMDRDNAarchitecture.
microbenchmark(seeFigure1),theresultssuggestedthatthefirst- • Wepresentvalidationresults,comparingNaviSimsimulation
resultsandhardwareexecutionontwoGPUs(AMDRX5500
levelcachelatencywasreducedfrom≈190cyclesto≈110cycles.
XTandAMDRadeonProW6800).
Thistypeofhiddenchangeareoftenignoredwhenconfiguring • WeprovidecasestudiesdemonstratingtheutilityofNaviSim,
thebaselineGPUimplementation,whichcanimpactthevalidity evaluatingnewarchitecturalanddesignfeaturestargeting
ofGPUarchitecturestudies. theRDNAGPU.
334

NaviSim:AHighlyAccurateGPUSimulatorforAMDRDNAGPUs PACT’22,October8–12,2022,Chicago,IL,USA
2 BACKGROUND wehaveobservedthatthework-groupdispatchingcanbeamajor
TheAMDGCNarchitecture[2]hasbeenthereferencespecifica- performancebottleneck.
tionforAMDGPUdesignoverthepastdecade.However,inrecent ACU(seeFigure3(c))onaGPUissimilartoaCPUcore.TheCU
years,theAMDGCNarchitecturefacedcriticalscalabilitychal- isresponsibleforinstructionexecutionanddataprocessing.Each
lenges,asitstruggledtofullyexploittheever-increasingnumberof CUincludesaschedulerthatcanfetchandissueinstructionsforup
transistorsprovidedbytoday’sphotolithographytechnology[8].In to40wavefronts.Duringeachcycle,theschedulercandecodeup
response,AMDdevelopedanewfamilyofRDNAarchitectures[3], to5differentinstructionsandissuethese5instructionstodifferent
introducingamajorarchitecturaloverhaulascomparedtotheGCN executionunits,includingabranchunit(notshowninthefigure
architecture.Inthissection,weprovideabriefoverviewoftheGCN forclarity),ascalarunit(responsibleforexecutinginstructionsthat
architectureanddiscussthechangesintroducedintheRDNAar- manipulatedatasharedbywork-itemsinawavefront),aLocalData
chitecture.Formoredetailedspecifications,readerscanrefertothe Share(LDS)unit,avectormemoryunit,andfourSingle-Instruction
AMDGCNWhitepaper[2]andRDNAWhitepaper[3]. Multiple-Data(SIMD)units.EachSIMDunitisresponsibleforex-
GPUProgrammingModel. BothAMDGCNandRDNAGPUs ecutingvectorizedfloating-pointinstructionsfor10outofthe40
canexecuteGPUprogramsimplementedusingthesameGPUpro- wavefrontsmanagedbythescheduler.EachSIMDunitisequipped
grammingmodel.Here,weintroducetheprogrammingmodelusing with16single-precisionArithmeticLogicUnits(ALUs).Therefore,
OpenCLterminologies,thoughotherprogrammingmodels(e.g., each64-work-itemwavefronttakes4cyclestofinishtheexecution
CUDAorHIP)usesimilarsemantics. ofoneinstruction.SinceCUbehaviorwilldetermineinstruction
AGPUprogramtypicallycanbeseparatedintotheGPUportion throughputandneedstohandledatadependencies,modelingand
andtheCPUportion.Theportionoftheprogramthatexecuteson analyzingthebehavioroftheCUaccuratelyisessentialtoaccu-
theGPUiscalledakernel.Theportionoftheprogramthatruns ratelymodelingtheoverallGPUperformance.
ontheCPU(i.e.,thehostprogram)launcheskernelsusingvendor- TheGCNarchitecturehasatwo-levelcachehierarchy.TheL1
providedruntimeAPIs.Akernelconsistsofanumberofwork- cachecanbedividedintotheL1scalarcache(mainlyusedforstor-
items.Work-itemsaresimilartoCPUthreads;theyexecutethe ingconstantdata,suchaskernelargumentsandpointers),anL1
sameprogramconcurrently,butworkondifferentdata.Work-items instructioncache,andanL1vectorcache(awrite-throughcache
canbegroupedintowork-groups.Thework-itemsinawork-group thatstoresmostofthedatarequiredbyaCU).EachCUhasaded-
(typically32–1024work-items)canbesynchronizedusingbarriers icatedL1vectorcache.CUsinaShaderArray(typically4CUs)
andcanshareasmall,butfast,localdatashare(LDS)memory. shareanL1scalarcacheandanL1instructioncache.AlltheL1
OnaGPU,asubsetofwork-items(typically32–64)inawork- cachesfetchdatafromL2caches(L2sarewrite-backcaches).Each
groupareorganizedintoawavefront.AMDGPUsissueinstructions L2cacheinterfacestoaDRAMcontroller(typicallyimplemented
atawavefrontgranularity—everytimetheinstructionscheduler inHBMorGDDRtechnology).TheL2cachesandtheDRAMcon-
issuesoneinstruction,theALUrepeatstheoperationforeachwork- trollersarebanked,allowingthemtoserviceapartoftheaddress
itemwithinthewavefront.Duetotheuseofthisformofwavefront- space.SincemanyGPUapplicationsarememorybound,careful
level scheduling, all the work-items in a wavefront are always modelingofthememorysystemiscriticalandagoalofNaviSim.
synchronized.Thismechanismisalsoknownaslock-stepexecution. TheL1cachesandtheL2cachesareconnectedwithacrossbar.
Withlock-stepexecution,allthework-itemsinawavefrontneedto Thecrossbardesigncanprovidelow-latencyandhigh-throughput
executeaninstruction,evenifonlypartofthewavefrontneedsit communicationchannelsfortheL1andtheL2caches.However,as
whentheexecutionpathsofdifferentwork-itemsdiverge.Thread the7nmtechnologyenablesalargernumberofCUs,thecrossbar
divergencenegativelyimpactsGPUperformance,asonlypartof designstrugglestoscale.Therefore,futurearchitectureswillrequire
thework-itemexecutionwilltakeeffect. designchangestargetingthememoryhierarchyinordertoenhance
GCN Architecture. The GCN architecture (see Figure 3(a)) thescalabilityoftheGPUs.
adoptsahighlymodulardesignthatincorporatesaCommandPro- RDNAArchitecture.AMD’sRDNAarchitecture(seeFigure3(b))
cessor,ShaderArrays(includingComputeUnitsandL1caches), isdesignedtoreplacetheGCNarchitectureforbetterscalability.
anon-chipnetworkconnectingthecore-sideL1cachesandthe TheRDNAarchitecturemakeschangestomanyelementsofthe
memory-sideL2caches,andDRAM.TheCommandProcessor(CP) GPUdesign,includingtheprogrammingmodel,theCU,andthe
isresponsibleforhandlingallcommunicationswiththeCPU,in- memoryhierarchy.
cludingmemorycopyingandkernellaunch.TheCPisalsorespon- OneofthemostnoteworthychangesisthattheRDNAarchi-
sibleforbreakingdownkernelsintowork-groupsandwavefronts, tecturereducesthesizeofthewavefrontfrom64work-itemsto
aswellasdispatchingthework-groupsandwavefrontstotheCom- 32work-items.Bycuttingthewavefrontsizeinhalf,theCUsare
puteUnits(CUs).IntheGCN3architecture,theCPcandispatch expectedtobettercopewithahigherdegreeofthreaddivergence
one wavefront per cycle. However, given the fact that GPU de- inmodernworkloads.Additionally,asthewavefrontsizeissmaller,
vicescanholdanever-increasingnumberofCUs,especiallywhen fewermemorytransactionsareexpectedtobegeneratedbyone
movingto7nmtechnology,theCPdispatchingspeedcanbeaper- load/storeinstruction(althoughthetotalnumberoftransactions
formancebottleneck.Whileothersimulators(e.g.,Multi2Sim[38], generatedbythekernelislikelytoremainunchanged),potentially
GPGPUSim[7])donottypicallymodelthecommunicationbetween reducingmemoryaccesslatencies(asdemonstratedinthepointer
theCPandtheCU,wecarefullymodelthedispatchingprocess,as chasingmicrobenchmarkresultsinFigure1).Thesebenefitsshould
helptoimproveALUutilization.Thispaperwillfurtherevaluate
theeffectsofnarrowerwavefrontsinSection6.
335

PACT’22,October8–12,2022,Chicago,IL,USA Bao,Sun,Feric,Shen,Weston,Abellán,Baruah,Kim,Joshi,andKaeli
Command Processor
Command Processor
|                     | Shader Array |     | Shader Array    |     |     |     |     |                   | Shader Array |     | Shader Array      |     |     |
| ------------------- | ------------ | --- | --------------- | --- | --- | --- | --- | ----------------- | ------------ | --- | ----------------- | --- | --- |
| L1I Instruction Mem |              |     |                 |     | L1I |     |     |                   |              |     |                   |     |     |
|                     |              |     | Instruction Mem |     |     |     |     | Dual Compute Unit |              |     | Dual Compute Unit |     |     |
Compute Unit
| L1V |              |     | Compute Unit |     | L1V |     |     |          |         |     |     |         |          |
| --- | ------------ | --- | ------------ | --- | --- | --- | --- | -------- | ------- | --- | --- | ------- | -------- |
|     |              |     |              |     |     |     |     | L0S      | L0V L0V | L0I | L0S | L0V L0V | L0I      |
| L1V | Compute Unit |     | Compute Unit |     |     |     |     | ehcaC 1L |         |     |     |         | ehcaC 1L |
L1V
L1V Compute Unit Compute Unit L1V Dual Compute Unit Dual Compute Unit
|     |              |     |              |     |     |     |     | L0S                | L0V L0V | L0I | L0S | L0V L0V            | L0I |
| --- | ------------ | --- | ------------ | --- | --- | --- | --- | ------------------ | ------- | --- | --- | ------------------ | --- |
| L1V | Compute Unit |     | Compute Unit |     | L1V |     |     |                    |         |     |     |                    |     |
|     |              |     |              |     |     |     |     |            ... ... |         |     |     |            ... ... |     |
L1S Scalar Mem Scalar Mem L1S 3-5 DCUs in Total 3-5 DCUs in Total
| Interconnect |     |     |     |     | Interconnect |     | Interconnect |     |     |     |     |     | Interconnect |
| ------------ | --- | --- | --- | --- | ------------ | --- | ------------ | --- | --- | --- | --- | --- | ------------ |
L2 Cache L2 Cache L2 Cache L2 Cache ... ... L2 Cache L2 Cache L2 Cache L2 Cache L2 Cache ... ... L2 Cache
|      |           |     |      | ... ... |      |     |      |     |      |      |      | ... ... |      |
| ---- | --------- | --- | ---- | ------- | ---- | --- | ---- | --- | ---- | ---- | ---- | ------- | ---- |
| DRAM | DRAM DRAM |     | DRAM |         | DRAM |     | DRAM |     | DRAM | DRAM | DRAM |         | DRAM |
(a) The Architecture of GCN3  GPUs.  (b) The Architecture of RDNA GPUs.
| L1I |              |     | L1S | L1V |     |     |     |     |     |         |     |     |     |
| --- | ------------ | --- | --- | --- | --- | --- | --- | --- | --- | ------- | --- | --- | --- |
|     | To Other CUs |     |     |     |     |     |     |     |     | L0I L0S |     |     | L0V |
Global
| 10 Wave  |                |     | Memory |            |     |      |           |       |       |       |            |           |      |
| -------- | -------------- | --- | ------ | ---------- | --- | ---- | --------- | ----- | ----- | ----- | ---------- | --------- | ---- |
| Inst Buf |                |     | Access | Register   |     |      |           |       | 20    | 20    |            |           |      |
|          |                |     |        | Write Back |     | SALU | SReg VReg | eussI | hcteF | hcteF | eussI VReg | SReg SALU |      |
|          | eussI & edoceD |     |        |            |     |      |           | Wave  |       | Wave  |            |           | Vec. |
10 Wave Vec.  LDS Operand Gathering Inst Inst Operand Gathering Mem
| hcteF tsnI |      |                   |      |      |      |          |         |     | Buf       | Buf       |       |          |          |
| ---------- | ---- | ----------------- | ---- | ---- | ---- | -------- | ------- | --- | --------- | --------- | ----- | -------- | -------- |
| Inst Buf   | VReg | SReg              | SALU | Mem  | Unit |          |         |     |           |           |       |          | tinU SDL |
|            |      |                   |      |      |      | SIMD-32  | Global  |     |           |           |       | SIMD-32  |          |
| 10 Wave    |      | Operand Gathering |      |      |      |          | Memory  |     |           |           |       |          |          |
| Inst Buf   |      |                   |      |      |      |  D-32    | Access  |     |           |           |       | SIMD-32  |          |
|            | SIMD | SIMD              | SIMD | SIMD |      |          |         |     | 20  hcteF | hcteF 20  | eussI |          |          |
10 Wave -16  -16  -16  -16  Operand Gathering eussI Wave Wave Operand Gathering Vec.
| Inst Buf |     |     |     |     |     |      |           |     | Inst | Inst |      |           | Mem |
| -------- | --- | --- | --- | --- | --- | ---- | --------- | --- | ---- | ---- | ---- | --------- | --- |
|          |     |     |     |     |     | SALU | SReg VReg |     |      | Buf  | VReg | SReg SALU |     |
Buf
(c) The Architecture of a GCN3 Compute Unit. (d) The Architecture of an RDNA Dual Compute Unit. L0V
Figure3:AComparisonbetweentheGCNArchitecture[2]andtheRDNAArchitecture[3].
AsecondmajorchangeinRDNAGPUsistheintroductionof (typically4-5)inaShaderArray,andsitsbetweentheL0andL2
DualComputeUnits(DCUs,seeFigure3(d)),replacingtheGCN caches.TheL1cachescanreducethenumberofrequestsarriving
CUs.ADCUcontains4schedulers.Increasingthenumberofsched- attheL2caches(inthecaseofL1hits)andreducetheamountof
ulersfrom1to4significantlyincreasestheinstructionissuerate. datathatistransmittedacrossthechip(fromL2toL0),therebyin-
Ratherthandispatchinginstructionsacross4SIMDunits,asinthe creasingperformanceandloweringthepowerconsumptioncaused
integratedGCNCU,eachRDNAschedulerinaDCUdispatches bycross-chiptransmissions.Finally,thecachelinesizeoftheL0
instructionsto1SIMDunit.OneSIMDunitinaDCUhas32single- vectorcaches,L1caches,andtheL2cachesisdoubledfrom64B
precisionALUs,doublingthenumberofALUsinaCU.Working to128B,sothatacachelinecandeliveruniquesingle-precision
numbersforall32work-itemsinawavefront(4𝐵×32=128𝐵).
togetherwiththenarrower32-work-itemwavefronts,eachSIMD
unitcanfinishexecutingoneinstructioninasinglecycle,ascom- The Akita Simulation Engine. The Akita Simulation En-
paredto4cyclesinaCU. gine[33]isacomputerarchitecturesimulatorenginethatisim-
Third,theRDNAarchitectureredefinesthecachehierarchyfrom plementedintheGoprogramminglanguage[37].TheAkitaSim-
a2-layerstructuretoa3-layerstructure.Thecachesthataredirectly ulation Engine has been used effectively in the MGPUSim sim-
connectedtotheDCUsarerenamedasL0caches(versusL1caches). ulator[33].WeselectedtheAkitaSimulationEnginebecauseof
Eachread-onlyL0instructioncacheandread-onlyL0scalarcache itshighflexibilityandoptimizedmulti-threadedsimulationper-
arenolongersharedbymultipleCUs,butdedicatedtoaDCU.Each formance.Whilewereuseafewcomponents(e.g.,theL0andL2
DCUconnectswith2separatewrite-throughL0vectorcaches;a caches)fromMGPUSim,mostoftheNaviSimsimulatorisdesigned
groupoftwoschedulersandtwoSIMDunitscanuseoneL0vector andimplementedindependentlyusingtheAkitaSimulationEngine
cache.SincenowwehavetwoL0cachesconnectedtooneDCU, (e.g.,theRDNAinstructionemulator,thenewL1write-evictcache,
updatingthedatainonecachemayrenderthedataintheother andDRAMcontrollers).
L0cachestale.ThismaycausecoherenceissueswithinaDCUand
requiresexplicitcacheinvalidationinstructions(asprovidedinthe
| RDNAISA,butnotintheGCNISA). |     |     |     |     |     |     | 3   | NAVISIM |     |     |     |     |     |
| --------------------------- | --- | --- | --- | --- | --- | --- | --- | ------- | --- | --- | --- | --- | --- |
Additionally,anintermediatelevelofcaching(i.e.,thenewwrite-
Inthissection,wepresentNaviSim,anovelGPUsimulatorthat
evictL1cache)isinserted.TheL1cacheservesagroupofDCUs modelstheAMDRDNAarchitecture.NaviSimisopensource(link
336

NaviSim:AHighlyAccurateGPUSimulatorforAMDRDNAGPUs PACT’22,October8–12,2022,Chicago,IL,USA
is hidden for double-blind review) under the terms of the MIT Scheduler
license[17].
RDNAISAEmulation. NaviSimisanexecution-drivensimu-
lator.ThesimulatorrecreatestheexecutionresultsofGPUinstruc-
tionsduringsimulationwiththehelpofaninstructionemulatorfor
theRDNAISA.NaviSimcanusetheMGPUSim’sGCN3instruction
emulatorasalibrary.Sincetheinstructionemulatorsharesthe
sameinterfacewiththeGCN3instructionemulator,thisallows
userstoswaptheinstructionemulatorbeingused.Thevirtualdri-
ver(asetofAPIsthatconnecttheGo-codedhostprogramsand
thesimulatedGPUs)ofNaviSimallowsuserstoconfigurewhich
ISAtoemulateandloadthecorrespondingGCN3/RDNAkernel
binaries.
NaviSimrunsineitheremulationmodeortimingsimulation
mode.Instructionemulationmodecanrecreateexecutionresults,
withoutevaluatingdetailedtiminginformationfortheinstruction
pipelines,caches,andDRAMcontrollers.Asaresult,emulationruns
muchfasterthantimingsimulation.Nomatterwhichmodeisused,
theemulationresultsofthebenchmarksshowninTable4exactly
matchtheoutputoftheapplicationsrunonrealGPUhardware.
Currently,NaviSimsupportsbothOpenCL[20]kernelsandker-
nelswrittenintheHIPprogramminglanguage[5].OpenCLkernels
canbecompiledbytheAMDofficialclang-oclcompiler,which
is a standard part of the AMD Radeon Open Compute (ROCm)
platform[35].HIPkernelscanbecompiledwiththehipcccom-
piler,whichalsoshipswiththeAMDROCmplatform.Byusing
the–gencoargument,hipccignoresthehostprogramandonly
generateskernelbinaries.NaviSimsupportsloadingkernelbinaries
compiledbyeithercompilerandemulatestheexecutionoftheGPU
kernelsusingactualinputdatasets.
WavefrontDispatching. WecarefullymodeltheCommand
Processortocapturethewavefrontdispatchingprocess.Themod-
eledCommandProcessormaintainsresourcemasksthatkeeptrack
ofwhichresourcesareoccupiedineachCU/DCU,includingthe
wavefrontslot(wavefront-levelresources,suchasinstructionbuffers
andtheprogramcounterregister),scalarregisters,vectorregisters,
andLDSmemory.Theresourcemasksensurethatnohardware
resourcesareoversubscribed.Sincewemasktheresourceseither
attheregisterlevel(forvectorandscalarregisters)oratthebyte
level(fortheLDS),wecanalsomodelregister/LDSfragmentation
issuesintheCUs/DCUs[29].
NaviSimcanalsosupportmodelingconcurrentkernelexecu-
tion[43],astheCommandProcessorhasmultiplewavefrontdis-
patchers.Eachwavefrontdispatcherwillmanagetheprogressof
the currently executing kernel’s execution and dispatch a new
wavefrontwhenresourcesfreeupintheCUs/DCUs.Bydefault,
eachCommandProcessorprovides8dispatchers(thenumberis
configurable),andhence,weallowupto8kernelstoexecutecon-
currently.Thewavefrontdispatcherscompetefortheresourcesand
haveequalopportunitiestodispatchwavefrontstotheCUs/DCUs,
ensuringthefairnessoftheconcurrentlyexecutingkernels.
DCUmodeling. IntheRDNAarchitecture,adualcomputeunit
replacestheoldcomputeunit.Wedevelopadetailedarchitectural
modeloftheDCU,whichgovernshowinstructionsareexecutedin
thesimulator.SincetheaccuratemodelingoftheDCUisessential
for simulation accuracy, NaviSim carefully models the pipeline
(seeFigure4)withamulti-stage,multi-issuestructure.
retibrA
hcteF
... ...
20 Wavefront
Inst Buffers in Total
... ...
retibrA
eussI
Decode ReadExec Write
Branch
Scalar
From 32-Width
L0 Inst ... ... ... SIMD
LDS
... Vector
Mem
50 Stages
To L0 Vector
Figure4:TheinstructionpipelinemodelinaNaviSimDCU.
EachDCUhasfourcopiesofthestructureshown.
Thepipelinestartswiththeinstructionfetcharbiter,whichat-
temptstofetchinstructionsforawavefrontthathasvacantspace
ininstructionbuffers.Whenmultiplewavefrontshavespaceavail-
able,thewavefrontthatreceivedinstructionsfurthestinthepast
isselected.Theissuearbitermonitorstheinstructionbuffersand
selectswavefrontsthathaveinstructionsreadythattheycanbe
issued.Thearbitercanissueatmost5instructionspercycle,one
instructiontoeachinstructionpipeline.
TheBranch,ScalarandLDSpipelinesuseafairlyregular6-stage
pipelinethatincludesfetch,issue,decode,read,execute,andwrite
stages.Notethatthedecodestagehappensaftertheissuestage.
Thisisbecausetheissuearbitercaneasilydeterminethetypeof
eachinstructionbycheckingafewencodedbitsintheinstruction.
Thebranchinstructionhasnodecodestage,astheinstructions
areverysimple.Thisisinlinewithpubliclyavailabledocumenta-
tionprovidedbyAMDfortheDCUarchitecture.TheSIMDunit,
whichprovidesmostofthecomputingpoweroftheDCU,also
usesa6-stagepipelinedesign,butiscapableofreading,executing,
andwriting32instructionsinparallel.Thevectormemoryhasthe
mostcomplexpipelinestructure,adoptinga50-stagepipeline.This
maysoundunusual,butthemodelmatchesourmicrobenchmark-
ingresults.Sincewearenotawareofthefunctionofeachofthe
pipelinestages,wedonotmodelthebehaviorofeachstage,but
onlyassociatealatencyvaluewitheachinstruction.
MemoryHierarchy. Wemodelthe3-levelcachehierarchyin
RDNAGPUs.L0scalarandL0instructioncachesareread-only,L0
vectorcachesuseawrite-throughpolicy,L1cachesuseawrite-
evictpolicy,andL2cachesuseawrite-backpolicy.Wealsoconnect
DRAMSim3[27]tomodelGDDR5,GDDR5X,GDDR6,GDDR6X,
HBM,andHBM2DRAMcontrollers.Weallowuserstoconfigure
anynumberofcachesinthehierarchyandallowanytypeofcombi-
nationsofthecachepolicies.Additionally,alloftheparametersof
acache(e.g.,setcount,wayassociativity,cachelinesize,directory
latency,storageaccesslatency)canbefullyconfiguredbyusers.
Multi-GPUandInterconnectModeling. SimilartoMGPUsim,
NaviSimnativelysupportsmulti-GPUsimulation.Thenumberof
GPUscanbeeasilyconfiguredwithcommandlinearguments(e.g.,
–gpus=1,2,3,4).ThereisnolimitationonthenumberofGPUs
undersimulation,aslongasthehostmachinehasenoughmemory.
Additionally,NaviSimnativelysupportsadvancedmulti-GPUfea-
tures,suchasunifiedmulti-GPUexecution[23],GPU-GPURDMA,
andGPU-GPUpagemigration[13].
UserInterface. NaviSimadoptsasimilaruserinterfaceasused
inMGPUSim.Torunbenchmarks,theusercancompileeitherHIP
337

PACT’22,October8–12,2022,Chicago,IL,USA Bao,Sun,Feric,Shen,Weston,Abellán,Baruah,Kim,Joshi,andKaeli
Table1:Theconfigurationsofvalidationplatforms. Table2:Simulatorconfiguration.
| Parameter | Platform1 | Platform2 |     | Param | P1 P2 | P3* |
| --------- | --------- | --------- | --- | ----- | ----- | --- |
GPU RadeonRX5500XT RadeonProW6800 BaseModel RX5500XT R9Nano -
| GPUCoreFreq  | 1845MHz   | 2320MHz   |     |          |                 |         |
| ------------ | --------- | --------- | --- | -------- | --------------- | ------- |
|              |           |           |     | #CU/DCU  | 11(DCU) 64(CU)  | 32(DCU) |
| DCUCount     | 11        | 30        |     | CoreFreq | 1845MHz 1000MHz | 1000MHz |
| GPUMemory    | GDDR6     | GDDR6     |     | TFLOPS   | 5.20 8.19       | 8.19    |
| MemBandwidth | 224.0GB/s | 512.0GB/s |     |          |                 |         |
|              |           |           |     | L0V$@    | 16KB 16KB       | 16KB    |
AMDRyzen
| CPU |     | AMDEPYC7302P |     | L0V$Assoc. | 4-way 4-way | 4-way |
| --- | --- | ------------ | --- | ---------- | ----------- | ----- |
Threadripper2950X
|           |                  |                  |     | L0Inst$@     | 32KB 32KB         | 32KB   |
| --------- | ---------------- | ---------------- | --- | ------------ | ----------------- | ------ |
| OS        | LinuxUbuntu18.04 | LinuxUbuntu20.04 |     | L0Scalar$@   |                   |        |
|           |                  |                  |     |              | 16KB 16KB         | 16KB   |
| GPUDriver | AMDROCm5.0       | AMDROCm5.1       |     |              |                   |        |
|           |                  |                  |     | L0I/L0S$Org. | PerDCU 4-CUshared | PerDCU |
|           |                  |                  |     | L1$          | 128KB -           | 128KB  |
|           |                  |                  |     | L1$Assoc.    | 16-way -          | 16-way |
orOpenCLkernelsusingtheofficialAMDcompiler.Ahostprogram L2$ 1MB 2MB 2MB
writteninGoisrequiredtoinvokeasetofAPIstoallocate/copy
|     |     |     |     | L2$Assoc. | 16-way 16-way | 16-way |
| --- | --- | --- | --- | --------- | ------------- | ------ |
memoryandlaunchkernels.ThehostAPIsarecompatiblewithMG-
PUSimandaresimilartocommonGPUprogrammingframeworks DRAMTech GDDR6 HBM HBM
(e.g.,CUDA,OpenCL).Userscanspecifyeitheremulationmode DRAMSize 4GB 4GB 4GB
|     |     |     |     | MemFreq | 1750MHz 500MHz | 500MHz |
| --- | --- | --- | --- | ------- | -------------- | ------ |
(valueemulationonly)ortimingmode(detailedtimingsimulation
|     |     |     |     | MemBus | 128bit 4096bit | 4096bit |
| --- | --- | --- | --- | ------ | -------------- | ------- |
thatinvolvesinstructionpipeline,caches,andmemorycontrollers)
|     |     |     |     | MemBW | 224GB/s 512GB/s | 512GB/s |
| --- | --- | --- | --- | ----- | --------------- | ------- |
asacommandlineoption.Meanwhile,NaviSimusesafewconfigu-
| rationfilesthatarewritteninGotodefinehardwareconfiguration. |     |     |     | *   |     |     |
| ----------------------------------------------------------- | --- | --- | --- | --- | --- | --- |
P3isnotanoff-the-shelfGPU,butitischosenspecificallyto
Userscaneasilyconfigurethehardwareundersimulationbymodi-
analyzetheperformanceimpactofindividualRDNAfeatures.
fyingparametersandcomponentconnectionsintheconfiguration @WeuseL0tonamethecachesthataredirectlyconnectedtoa
code.Weusecodeforconfigurationbecauseuserscaneasilydebug CU/DCU.InP2,L0cachesconnectstoL2cachesdirectly.
theconfigurationlogicwithdebuggers.
| SimulatorOutput. | NaviSimcangenerateawiderangeofout- |     |     |     |     |     |
| ---------------- | ---------------------------------- | --- | --- | --- | --- | --- |
putdatatofacilitateperformanceanalysis.Forhigh-levelmetrics,
NaviSimoutputsthetotalexecutiontime(kerneltime+memory PlatformswithGPUDevices. WeuseanRX5500XTGPU
copytime),totalkernelexecutiontime,andtheper-GPUkernel andRadeonProW6800(seeTable1)tovalidateourNaviSimGPU
executiontimes.Forperformancemetricsrelatedtoindividualcom- Model.OurplatformsruntheROCm5.0/5.1softwarestackson
ponents,NaviSimreportsinstructioncounts,averagelatencyspent
LinuxUbuntu18.04/20.04servers.
accessingeachlevelofcache,transactioncountsforeachcache(in-
SimulatorConfiguration. WesetthebaselineGPUconfigu-
cludingreadmisses,readhits,readMSHRhits,writemisses,write
rationusingpubliclyavailableinformationandcalibrateourresults
hits,andwriteMSHRhits),TLBtransactioncounts(hits,misses, usingmicrobenchmarks.ThedefaultconfigurationofNaviSimfor
MSHRhits),DRAMtransactioncountsandread/writesizes,and theRX5500XTGPUsisrecordedasP1inTable2.
transactioncountsfortheGPURDMAengines.
AftervalidatingtheNaviSimGPUmodel,weusethemodelto
| Additionally, | NaviSim | can generate low-level | detailed traces, |     |     |     |
| ------------- | ------- | ---------------------- | ---------------- | --- | --- | --- |
conductaseriesofuse-caseexperimentstoevaluatetheperfor-
includinginstructiontraces(completewiththeregisterstatesaf-
manceimpactofmicroarchitecturaldesignfeaturesintheRDNA
terexecutingeachinstruction)andmemorytraces(ateachlevel architecture.Wecompareoursimulationresultswiththedefault
ofcacheandtheDRAM,includingthetransactionstartandend configurationofMGPUSimfortheR9NanoGPUs(denotedasP2in
times).NaviSimshipswithagraphicaluserinterface(GUI)tool
Table2).SincetheAMDRX5500XTandtheR9NanoGPUsbelong
thatallowsuserstonavigatetheinstructiontracesandinspecthow
totwoverydifferentmarkets,comparingthemdirectlywouldnot
theregistersareupdatedaftereachinstructionexecution.Thistool
provideuswithalotofnewinsights.Therefore,wehaveconfig-
issimilartopopularGUI-basedMIPSemulators[28]andcanbe uredaGPUasP3inTable2.WechosetheseGPUconfigurationsso
usedforeducationalpurposes.Finally,NaviSimcanproducetraces thattheywouldhaveasimilartheoreticalcomputingthroughput
usingtheDaisenformat[36],sothatuserscanuseavisualization
(representedbyTFLOPS)andacomparablememoryhierarchy.In
tooltoinspectthedetailedbehaviorofeachcomponent.
theexperiments,ratherthandirectlycomparingP2andP3,we
graduallyaddfeaturesfromP2tobuildP3,sothatwecanevaluate
4 METHODOLOGY
theeffectofeachfeature.
Thestudiestobeperformedinthispaperincludetwoparts.First, Microbenchmarks. Weusemicrobenchmarkstoevaluatekey
wecalibrateandvalidatetheaccuracyofthesimulatorwithasetof parametersofeachGPUinordertoevaluatetheaccuracyofNavi-
microbenchmarksandfullbenchmarks.Withavalidatedsimulator, Sim.Wedesign7microbenchmarks(seeTable3)toevaluateindivid-
we then conduct additional experiments to evaluate the RDNA ualGPUsubsystems,suchasthewavefrontdispatcher,instruction
architecturedesign. pipelines,andmemoryhierarchy.Atahighlevel,werepeatone
338

NaviSim:AHighlyAccurateGPUSimulatorforAMDRDNAGPUs PACT’22,October8–12,2022,Chicago,IL,USA
Table3:ThelistofmicrobenchmarksusedtocalibrateNaviSim.
Microbenchmark Parameter Description
empty_kernel Work-Group Executesakernel,whilevaryingthekernelandwork-groupsizes.Weleave
DispatchingSpeed thekernelblanksothatwecanmeasurework-groupdispatchinglatency.
empty_kernel_multi KernelLaunching Launchesagivennumberofemptykernels.Thekernelusedissameastheone
Overhead usedinempty_kernelmicrobenchmark.
single_thread_loop InstructionPipeline Executesakernelwithonlyonethread.Thekernelhasamainloop,andeach
Depth iterationoftheloopexecutesasinglesingle-precisionfloatingpointinstruction.
multi_thread_loop InstructionPipeline Executesthesamekernelassingle_thread_loop,butwithalargenumberof
Throughput threadsthatarelargeenoughtofullyoccupytheGPUexecutionresources.
pointer_chasing_random CacheSizesand Runstheclassicpointerchasingmicrobenchmark[41].Thepointersstoredin
CacheLatencies thebuffersarerandomized.Weonlyrun1threadinthiskernel.
pointer_chasing_linear CacheSizeand Thesamekernelaspointer_chasing_random,exceptthatthepointersinthe
CacheLatencies buffersalwayspointtothenextpointer(thelastpointerpointstothefirstone.
memory_copy Bandwidthof Copiesdatafromonebuffertoanother,usingagivennumberofwork-groups.
CachesandDRAMs Whenthereareonly1-2work-groups,weevaluatetheL0cachebandwidth.
Whenthereareenoughwork-groupstofillthewholeGPU,thismicrobench-
markcanalsotesttheDRAMbandwidth.
operationthousandstomillionsoftimestostressindividualcom- Table4:FullBenchmarks.
ponentsoftheGPU.Thenwestatisticallyanalyzethelatencyofan
operationorevaluatethethroughputofaspecificGPUcomponent. Abbr. Suite Workload
Asanexampleofthisprocess,weusethepointer-chasingbench-
ATAX PolyBench MatrixTransposeandVectorMulti-
marktoevaluatecachesizesandlatencies.Thehostprogramfirst
plication
createsaregionofmemory(i.e.,thepointerchasingregion)ofa
BICG PolyBench BiCGStabLinearSolver[39]
givensize.Next,thehostprogramdividestheregioninto8-byte
BS AMDAPPSDK BitonicSort
cellsandfillseachcellwithanaddressthatpointstothenextcell,
FIR HeteroMark FiniteImpulseResponseFilter
locatedatarandomlyassignedaddresswithintheregion.There
FLW AMDAPPSDK Floyd-WarshallAlgorithm
isnorepetitionintheaddressesstoredintheregionsothatthe
FWT AMDAPPSDK FastWalshTransform
wholeregioncanbetraversedmultipletimes.TheGPUwillusea
KM HeteroMark KMeansClustering
singlethreadtoaccessthecellsmanytimes(atleastseveraltimes
MT AMDAPPSDK MatrixTranspose
morethanthenumberofcells),followingtheaddressesstoredin
ReLU DNNMark RectifiedLinearUnit
theregion.Weeventuallydividethekernelexecutiontimebythe
SPMV SHOC SparseMatrix-VectorMultiplication
numberofaccessestocalculatetheaverageaccesslatency.
FullBenchmarks. WealsoexploitalistoffullGPUbench-
marks(seeTable4)fromawiderangeofbenchmarksuitesincluding
executiontimeover10runs.ThetimesobtainedontherealGPU
AMDAPPSDK[31],SHOC[10],HeteroMark[34],PolyBench[30],
arerecordedusingOpenCLevents[20].
andDNNMark[11].Weusethesebenchmarksasasetofwork-
Weusekernelexecutiontimeastheprimarymetrictoevaluate
loadsthatcoverawiderangeofapplicationstocomprisedifferent
theaccuracyofNaviSim,fortworeasons.First,executiontimeis
arithmeticintensities,memoryaccesspatterns,andcommunication
themostcommonlyreportedmetricwhenconsideringarchitectural
patterns.WestartwiththeoriginalOpenCLkernelimplementa-
tradeoffs.Second,thekernelexecutiontimeisahigh-levelmetric
tionsfromthebenchmarksuiteandcompilethekernelswiththe
thatsummarizestheimpactofallfeaturesbeingsimulated;wecan
originalAMDROCmcompiler(applyingdefaultcompileroptimiza-
onlyachievealowerrorinexecutiontimeifthefidelityofallthe
tions).WealsowritehostprogramsinGo,allowingthesimulators
componentsbeingmodeledinthesimulatorishigh.
tocallthekernels.Weensureourhostprogramsareequivalentto
theoriginalhostprogramfromthebenchmarksuite.Forvalidation 5 SIMULATORVALIDATION
experiments,wevarytheproblemsizesofthebenchmarkstomake
Anycomputerarchitecturesimulatorrequiresarigorousvalidation
surethatNaviSimcanrecreatescalingtrends.Forperformance
processbeforeitcanserveasabaselineforfutureresearch.Inthis
evaluation,weuselargeproblemsizesthataresufficienttostress
section,wereportonourvalidationeffortsforNaviSim,comparing
thewholeGPU.
simulationresultsagainstGPUhardwareexecution.
RunningBenchmarks. WerunOpenCLimplementationsof
First,weverifythecorrectness(intermsofapplicationoutputs)
thebenchmarks,whilevaryingtheinputsize.Weusekernelex-
ofNaviSiminbothemulationmodeandtimingsimulationmode.
ecutiontimeastheperformancemetricandreporttheaverage
Tothisaim,wecompareeverysimulator-generatedapplication
outputwithitscorrespondingactualhardwareexecutionoutput.
339

PACT’22,October8–12,2022,Chicago,IL,USA Bao,Sun,Feric,Shen,Weston,Abellán,Baruah,Kim,Joshi,andKaeli
800
400
200
100
60
32 128 512 2K 8K 32K 128K 512K 2M 8M 32M
Pointer Chasing Region Size (Byte)
)sn(
sseccA
reP
emiT
NaviSim RX 5500 XT GPU executionwithhighfidelity,weclearlycapturethepatternsandthe
nuancesoftheGPUarchitecture.Forexample,intheFIRandReLU
benchmarks,wesuccessfullycapturethetransitionintheworkload
whentheexecutiontimestartstoincrease.Thissuggeststhatwe
properlymodelanylimitationsassociatedwiththetotalamountof
computingresourcesavailableonaGPU.Asanotherexample,we
observestepsinsomebenchmarks(e.g.,BSbefore32K,FWTafter
16Kand32K,KMbetween2Kand3K).Asthesestepsarecausedby
complexinteractionsbetweentheinstructionscheduler,cachehier-
archy,andmemorytransactionhandling,beingabletomodelthese
Figure 5: Comparison of the Pointer Chasing microbench- stepsdemonstratesthatNaviSimcanmodelthesubtlefeaturesin
markbetweenNaviSimandtheRX5500XThardware. theRDNAarchitecture.WehavealsovalidatedNaviSimagainsta
secondGPUmodel,theAMDRadeonProW6800GPU(seeTable1),
whichisanRDNA2-basedGPU.Theaveragedifferencebetween
Weobservethattheresultsmatchexactlyforallapplications.The simulatedandhardwaremeasuredexecutiontimeis10.08%.Our
matchingresultsofNaviSimwhenrunningintheemulationmode validationresults,whichareshowninFigure7,demonstratethat
suggeststhatNaviSimfaithfullyemulatestheRDNAinstructions. NaviSimcancapturethechangesofmicroarchitecturesinRDNA
Additionally,differentfromothersimulators[7,38],wherethein- acrossdifferentdeviceswithhighfidelity.
structionemulationandtimemodelingareindependentofeach Wearealsoawareofdiscrepanciesinafewbenchmarks,suchas
other,NaviSimmodelsdatavaluesineveryrequestandeverysim- FWandKM.Ingeneral,thesebenchmarksareeithershort-running
ulatedclockcycle.Oursimulationapproachexposesmistakesin benchmarks (FW) or workloads that involve a large number of
thecommunicationmodeling,manytimescapturedasanerror kernellaunches(KM).ThediscrepanciessuggestthatNaviSimhas
inthesimulator-generatedexecutionoutput.Beingabletomatch difficultiesinmodelingGPUbehaviorsatthekernellaunchphase.In
theresultsbetweenthetimingsimulationresultsandhardware general,webelievethisisnotabigproblemsincethesimulatorsare
measurementssuggeststhatwefaithfullymodelcommunicationbe- likelytobeusedtomodellargeproblemsizes,andthedifferences
tweencomponents.Forexample,wecaptureandaccuratelymodel observed in the kernel launch overhead should not impact the
allflushesofcachelines,whereerrorsmayneverbedetectedby overallaccuracybymuch.Weleavemoredetailedmodelingofthe
simulatorsthatseparatetheemulationandtimingmodelinglogic. kernellaunchbehaviorasfuturework.
Blending the instruction and emulation code provides an extra WealsoanalyzethememoryfootprintandperformanceofNavi-
layerofconfidencethatthecomponentsandthecommunication Simsimulation.IntheFIRbenchmark(witha4Mproblemsize,as
aremodeledaccurately. showninFigure6),weuse823MBofmemory,whichfitsinthemem-
Second,forvalidatingthetimingmodel,weusebothmicrobench- oryofmostmoderncomputers.Intermsofsimulationperformance,
marksandfullapplicationworkloadstovalidatetheaccuracy.We on an Apple M1 Mac Mini, we achieve 43.5KIPS and 89.5KIPS
configureoursimulationusingthepubliclyavailableparametersof inserialandparallelmodes(NaviSimismulti-threaded),respec-
theAMDRX5500XTGPU(seeTable1).Wethenusemicrobench- tively. This performance is much faster than MGPUSim, which
marks(seeTable3)tohelpusreverse-engineerawiderangeof reported27KIPSparallelexecutionperformanceintheiroriginal
undocumentedparameters,suchasthework-groupdispatchover- paper.Thus,thememoryconsumptionandperformanceofNaviSim
head,instructionpipelinedepth,cachelatencyateachlevel,and arequitereasonable.
DRAMbandwidthandlatency.Agoodexampleofourstrategyis
thatweuseapointerchasingmicrobenchmarktofigureoutthe 6 CASESTUDIES:UNDERSTANDINGTHE
sizeandlatencyofeachlevelofthecache.Aswecanseefromthe RDNAARCHITECTUREFEATURES
representativeresultsofrunningthepointerchasingmicrobench-
Withacarefullyvalidatedsimulatormodel,wenextuseNaviSim
mark(seeFigure5),NaviSimisabletocalibratetheparameters
toperformasetofexperimentstoanalyzetheimpactofRDNA
withextremelyhighaccuracy.Thecalibrationresultsofothermi-
featuresonapplicationperformance.Inparticular,weattemptto
crobenchmarksfollowsimilartrendsasthepointerchasingresults.
answerthefollowingquestions:
Withthefullycalibratedparameters,weevaluatetheaccuracy
ofthesimulatorbyvalidatingitagainstAMDRX5500XThard- (1) HowdoestheISAimpacttheoverallperformanceandhow
wareusingfullbenchmarks(Table4).Figure6showsthesimulator does the DCU architecture impact performance? How is
accuracy(lefty-axis)aswellastherelativeerrorateachdatapoint theinstructionexecutionpipelineimpactedbyexecutinga
bythebarplot(righty-axis)forarangeofproblemsizesdepending differentISA?
uponeachbenchmark.Weobservethatforsomebenchmarksthe (2) WhatistheeffectofthenewlyaddedL1cache?
relativeerrorincreasesastheproblemsizeincreases(e.g.,ReLU), (3) Whatpercentoftheoverallperformanceincreasecanbe
and for some other benchmarks the relative error decreases as attributedtochangesinfrequency(increasedfrom1GHzto
theproblemsizeincreases(e.g.,BS,SPMV).Foreachbenchmark, 1.845GHz)?
inparenthesis,wealsoreporttheaverageerror.Aswecansee, Whilewestudytheimpactofthefeaturesontheoverallperfor-
theerrorintermsofmodeledexecutiontime,averagedacrossall mance,wealsousethisstudytodemonstratetheutilityofNaviSim.
benchmarks, is just 9.75%. Additionally, we not only model the WeshowcasetheflexibilityandtheconfigurabilityofNaviSim.
340

NaviSim:AHighlyAccurateGPUSimulatorforAMDRDNAGPUs PACT’22,October8–12,2022,Chicago,IL,USA
|     |     | Simulation |     | Native Execution |     | Relative Error |     |
| --- | --- | ---------- | --- | ---------------- | --- | -------------- | --- |
)%( rorrE evitaleR
60
| ATAX (7.9%) |     | BICG (5.9%) |     |     |     | 12 FIR (12.9%) | 16  |
| ----------- | --- | ----------- | --- | --- | --- | -------------- | --- |
| 300         |     | 300         |     | 400 |     |                |     |
| )su(emiT    |     | 12          |     | 16  |     |                |     |
|             |     |             |     | 300 |     | 9 40           | 12  |
| 200         |     | 9 200       |     | 12  |     |                |     |
6
|     |     |     |     | 8 200 |           |     | 8   |
| --- | --- | --- | --- | ----- | --------- | --- | --- |
|     |     | 6   |     |       |           | 20  |     |
| 100 |     | 100 |     |       |           | 3   | 4   |
|     |     | 3   |     | 4 100 | BS (3.9%) |     |     |
| 0   |     | 0 0 |     | 0 0   |           | 0 0 | 0   |
64 512 1K 1.5K 2K 64 512 1K 1.5K 2K 4K 16K 32K 48K 64K 64 1K 16K 256K
|            |     |     |     | 120 |     | 300           | )%( rorrE evitaleR |
| ---------- | --- | --- | --- | --- | --- | ------------- | ------------------ |
| 8          |     |     |     | 8   |     |               |                    |
| FW (19.0%) |     | 60  |     |     |     | 16 MT (13.5%) |                    |
| )su(emiT   |     | 32  |     | 90  |     |               | 32                 |
| 6          |     |     |     | 6   |     | 12 200        |                    |
24
|     |     | 40           |     | 60  |            |     | 24  |
| --- | --- | ------------ | --- | --- | ---------- | --- | --- |
| 4   |     |              |     | 4   |            | 8   |     |
|     |     | 16           |     |     |            | 100 | 16  |
| 2   |     | 20           |     | 30  |            |     |     |
|     |     | 8 FWT (5.5%) |     | 2   | KM (11.5%) | 4   | 8   |
| 0   |     | 0 0          |     | 0 0 |            | 0 0 | 0   |
16 64 128 192 256 4K 16K 32K 48K 64K 256 1K 2K 3K 4K 256 1K 1.5K 2K
| 60  |     | 16  |     | 40  |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- |
100
| ReLU (3.1%) |     | SPMV (9.8%) |     |     |     |     |     |
| ----------- | --- | ----------- | --- | --- | --- | --- | --- |
| )su(emiT    |     | 12          |     | 30  |     |     |     |
| 40          |     | 75          |     |     |     |     |     |
|             |     | 8           |     | 20  |     |     |     |
50
20
|     |     | 4   |     | 10  |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- |
25
| 0      |            | 0 0    |         | 0   |     |     |     |
| ------ | ---------- | ------ | ------- | --- | --- | --- | --- |
| 64 512 | 8K 128K 1M | 64 512 | 1K 1.5K |     |     |     |     |
Figure6:SimulatorvalidationagainsttheAMDRX5500XTGPU.Thex-axisplotstheproblemsizeandthetwoy-axesplotthe
kernelexecutiontimeandrelativeerror.ThenumbersintheparenthesesrepresenttheaveragedifferencebetweenNaviSim
simulationandthehardwareexecutionforeachapplication.
ImpactofchangesintheISAandtheintroductionofthe vectorL0caches,thebandwidthbetweenthecoreandtheL0caches
DCU. Inthefirstsetofexperiments,wefocusonquestion(1).We iseffectivelydoubled,causinganincreaseintheperformanceby
| usetheP2configuration(seeTable2,denotedasCU+GFX803in |     |     |     | about2×. |     |     |     |
| ---------------------------------------------------- | --- | --- | --- | -------- | --- | --- | --- |
Figure8)asthebaseline.Wehaveconfigurationswhereweeither ThemostcriticaldifferenceintheISAsisthewavefrontsize
changetheISAtoGFX1010(CU+GFX1010)orchangethecoreto difference.KernelscompiledtotheGFX1010ISAalwaysuseawave-
DCU (DCU+GFX803, where the number of DCUs is half of the frontsizeof32,whichisnotaperfectmatchfortheCUarchitecture.
numberintheCU,thusprovidingafaircomparisonofsimilarcom- WhenrunningwavefrontsonCUs,becausetheCUschedulercan
putingcapabilities).Wealsoprovideaconfigurationthatchanges onlyissue1instructiontoaSIMDunitevery4cycles,theCUneeds
boththecoreandtheISA(DCU+GFX1010). tospend4cyclestoexecute1instruction.Weobserveunderutiliza-
At a high level, our results (see Figure 8) suggest that many tionoftheALUsinATAX,BICG,FIR,andFLWbenchmarkswhen
benchmarks(e.g.,BS,PR,ReLU)achievethesameperformancefor wematchtheGFX1010ISAwiththeCUmicroarchitecture.Other
allfourcases.Thisisunderstandable,asthememorysystemand benchmarksdonotobservethisissuebecausetheyareboundby
theoverallcomputingcapabilitiesremainthesame.However,we memorybandwidthandarenotsensitivetoALUutilization.Inthe
noticemajorperformancedifferencesinATAX,BICG,FIR,andFLW FIRbenchmark,wealsoseethatmatchingtheGFX803benchmark
benchmarks.Thesechangesarecausedbythedifferencesinthe andtheDCUmicroarchitecturecausestheSIMDunitnottobeable
ISAandtheCU/DCUorganization. tocatchupwiththeinstructionissuingspeed,leadingtosignificant
ATAXandBICGrevealhowtheCU/DCUorganizerimpactsthe pipelinestallingandevenmoreslowdown.Theproblemdoesnot
performanceofbenchmarksthataresensitivetomemoryband- existwhentheDCUmicroarchitectureandGFX1010ISAareused
widthandlatency.ATAXandBICGareworkloadswithlimited together,suggestingthattheISAandthemicroarchitectureare
parallelismandstronginter-work-itemdependencies.Therefore, co-designedtoachievethebestperformance.
ATAXandBICGhavelargework-groupsizesandasmallnumberof ImpactoftheL1cache. Next,wefocusonquestion(2)and
work-groupsineachkernel.IntheCUconfigurations,thenumber trytounderstandtheeffectofthenewL1cache.WeusetheP3
ofblockscannotutilizealltheCUsintheGPUandhence,cannot configuration(seeTable1,denotedasP3w/L1inFigure9)asthe
fullyutilizethebandwidthbetweentheL0cachesandtheCUs. baseline,andcomparetheperformancewithaconfigurationthat
Onthecontrary,thenumberofblocksintheATAXandBICGcan removestheL1cache(P3w/oL1).TheresultsareshowninFigure9.
fullyutilizetheDCUsonthethirdandfourthconfigurations,since Similartotheearlierexperiments,asweonlymodifyasmall
thecountofDCUsishalved.AseachDCUisconnectedwithtwo partoftheconfiguration,mostbenchmarksthatareboundedby
341

PACT’22,October8–12,2022,Chicago,IL,USA Bao,Sun,Feric,Shen,Weston,Abellán,Baruah,Kim,Joshi,andKaeli
|     |     |     |     |     | Simulation |     | Native Execution |     |     | Relative Error |     |     |     |
| --- | --- | --- | --- | --- | ---------- | --- | ---------------- | --- | --- | -------------- | --- | --- | --- |
)%( rorrE evitaleR
|          | 300         |     |     |     |                  |     | 1200 |     |     | 8   | 20          |     |     |
| -------- | ----------- | --- | --- | --- | ---------------- | --- | ---- | --- | --- | --- | ----------- | --- | --- |
|          | ATAX (9.3%) |     |     | 30  | 300 BICG (11.4%) |     |      |     |     |     | FIR (11.0%) |     | 30  |
| )su(emiT |             |     |     |     |                  |     | 120  |     |     | 6   | 15          |     |     |
|          | 200         |     |     |     |                  |     | 800  |     |     |     |             |     |     |
|          |             |     |     | 20  | 200              |     |      |     |     |     |             |     | 20  |
|          |             |     |     |     |                  |     | 80   |     |     | 4   | 10          |     |     |
|          | 100         |     |     |     |                  |     | 400  |     |     |     |             |     | 10  |
|          |             |     |     | 10  | 100              |     | 40   |     |     | 2   | 5           |     |     |
BS (3.3%)
|     | 0   |     |     | 0   | 0   |     | 0   | 0   |     | 0   | 0   |     | 0   |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
64 512 1K 1.5K 2K 64 512 1K 1.5K 2K 4K 16K 32K 48K 64K 64 1K 16K 256K
|     | 8        |     |     |     |     |     | 8 160 |     |     |     | 60         |     | 40 )%( rorrE evitaleR |
| --- | -------- | --- | --- | --- | --- | --- | ----- | --- | --- | --- | ---------- | --- | --------------------- |
|     |          |     |     | 32  |     |     |       |     |     | 16  | MT (14.6%) |     |                       |
|     | )su(emiT |     |     |     | 90  |     | 6 120 |     |     |     |            |     | 30                    |
|     |          |     |     |     |     |     |       |     |     | 12  | 40         |     |                       |
24
|     | 4   |     |     |     | 60  |     | 4   | 80  |     |     |     |     | 20  |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
|     |     |     |     | 16  |     |     |     |     |     | 8   |     |     |     |
20
|     |     |           |     |     | 30         |     | 2   | 40        |     |     |     |     | 10  |
| --- | --- | --------- | --- | --- | ---------- | --- | --- | --------- | --- | --- | --- | --- | --- |
|     |     | FW (8.9%) |     | 8   | FWT (4.2%) |     |     | KM (7.0%) |     | 4   |     |     |     |
|     | 0   |           |     | 0   | 0          |     | 0   | 0         |     | 0   | 0   |     | 0   |
16 64 128 192 256 4K 16K 32K 48K 64K 256 1K 2K 3K 4K 256 1K 1.5K 2K
|     | 30          |     |     | 15  | 300          |     | 100 |     |     |     |     |     |     |
| --- | ----------- | --- | --- | --- | ------------ | --- | --- | --- | --- | --- | --- | --- | --- |
|     | ReLU (6.3%) |     |     |     | SPMV (24.7%) |     |     |     |     |     |     |     |     |
)su(emiT
|     | 20  |     |     | 10  | 200 |     |     |     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
60
|     | 10  |     |     | 5   | 100 |     |     |     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
20
0
|     | 0   |     |         | 0   |        |     | 0    |     |     |     |     |     |     |
| --- | --- | --- | ------- | --- | ------ | --- | ---- | --- | --- | --- | --- | --- | --- |
|     | 64  | 512 | 8K 128K | 1M  | 64 512 | 1K  | 1.5K |     |     |     |     |     |     |
Figure7:SimulatorvalidationagainsttheAMDRadeonProW6800GPU.Thex-axisplotstheproblemsizeandthetwoy-axes
plotthekernelexecutiontimeandrelativeerror.Thenumbersintheparenthesesrepresenttheaveragedifferencebetween
NaviSimsimulationandthehardwareexecutionforeachapplication.
|                |     |     | CU+GFX803  |     | DCU+GFX803  |     |     |          |     | P3 w/o L1 | P3 w/ L1 |     |     |
| -------------- | --- | --- | ---------- | --- | ----------- | --- | --- | -------- | --- | --------- | -------- | --- | --- |
|                |     |     |            |     |             |     |     | 1.0      |     | 1.0       |          | 1.0 |     |
| emiT noitucexE |     |     | CU+GFX1010 |     | DCU+GFX1010 |     |     | 0.8      |     | 0.8       |          | 0.8 |     |
| dezilamroN     |     |     |            |     |             |     |     | etar tiH |     |           |          |     |     |
|                | 1.5 |     |            |     |             |     |     | 0.6      |     | 0.6       |          | 0.6 |     |
|                |     |     |            |     |             |     |     | 0.4      |     | 0.4       |          | 0.4 |     |
1.0
|     |          |     |        |         |       |             |     | 0.2                         |                 | 0.2       |                              | 0.2            |                                   |
| --- | -------- | --- | ------ | ------- | ----- | ----------- | --- | --------------------------- | --------------- | --------- | ---------------------------- | -------------- | --------------------------------- |
|     | 0.5      |     |        |         |       |             |     | 0.0                         |                 | 0.0       |                              | 0.0            |                                   |
|     |          |     |        |         |       |             |     | XATA GCIB SB RIF WLF TWF MK | TM RP ULeR VMPS | XATA GCIB | SB RIF WLF TWF MK TM RP ULeR | VMPS XATA GCIB | SB RIF WLF TWF MK TM RP ULeR VMPS |
|     | ATAXBICG |     | BS FIR | FLW FWT | KM MT | PR ReLUSPMV |     |                             |                 |           |                              |                |                                   |
|     |          |     |        |         |       |             |     | L0 Cache                    |                 |           | L1 Cache                     |                | L2 Cache                          |
Benchmark
Figure10:Thecachehitrateofeachlevelsofcachesbefore
Figure8:TheeffectofchangingfromCUstoDCUsandthe
andaftertheL1cacheisadded.
ISAfromGFX803(GCN3)toGFX1010(RDNA)onexecution
time.
|     |     |     |     |     |     |     |     |     | P3@1000MHz |     | P3@1845MHz |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | ---------- | --- | ---------- | --- | --- |
2.0
1.5
pudeepS 1.5
| emiT noitucexE |     |     |     | P3 w/o L1 | P3 w/ L1 |     |     |          |        |     |           |             |          |
| -------------- | --- | --- | --- | --------- | -------- | --- | --- | -------- | ------ | --- | --------- | ----------- | -------- |
| dezilamroN     | 1.0 |     |     |           |          |     |     | 1.0      |        |     |           |             |          |
|                | 0.5 |     |     |           |          |     |     | 0.0      |        |     |           |             |          |
|                |     |     |     |           |          |     |     | ATAXBICG | BS FIR | FLW | FWT KM MT | PR ReLUSPMV | Geo-Mean |
Benchmark
0.0
|     | ATAXBICG |     | BS FIR | FLW FWT   | KM MT | PR ReLUSPMV |     |                                                    |     |     |     |     |     |
| --- | -------- | --- | ------ | --------- | ----- | ----------- | --- | -------------------------------------------------- | --- | --- | --- | --- | --- |
|     |          |     |        | Benchmark |       |             |     | Figure11:Theeffectoftheincreasedcorefrequencyonthe |     |     |     |     |     |
overallperformance.Asweincreasethecorefrequencyby
Figure9:Theeffectofaddingthenewwrite-evictL1cache
1.845×,theperformanceimprovedby≈1.5×.
onbenchmarks’executiontime.
TounderstandwhyFIRandFLWcanbenefitfromtheaddedL1
either compute power or by DRAM bandwidth do not observe cache,weplotthecachehitrateforeachlevelofcacheinFigure10.
aperformancedifference.However,wedidobserveperformance Ingeneral,weseeahighdegreeofdiversityacrossthebenchmarks,
improvementsintheFIRandFLWbenchmarks. astheseworkloadshavedramaticallydifferentlevelsoflocality.For
342

NaviSim:AHighlyAccurateGPUSimulatorforAMDRDNAGPUs PACT’22,October8–12,2022,Chicago,IL,USA
theBS,FIR,FLWandFWTbenchmarks,weobservearelatively thepotentialbenefitsofeacharchitecturalfeature.Forexample,
highL1cachehitrate(around40%).However,theeffectsofthehigh wecanchangetheCUswithDCUs,withoutimpactinganyother
L1cachehitontheL2cachesdifferacrossworkloads.ForBSand features,includingtheGPUfrequency.Thisisverychallengingto
FWT,theL2cachehitrateremainsunchanged,whiletheL2cache dowithRTLandimpossibletodoonlivehardware.
hitratedropssignificantlyforFIRandFLW(inagoodway).These However,wealsofindsomechallengeswhendesigningthese
resultssuggestthattheimprovedL1cachehitratescanreducethe experiments.Onemajorchallengeisdeterminingthevalidityof
numberoftransactionstotheL2cachesfortheFIRandtheFLW themodeledconfigurations.Forexample,earliersimulators(e.g.,
benchmarks.BecausetheFIRbenchmarkalreadyhasahighL0hit MGPUSim) simply use an ideal network to connect the L1 and
rate,thespeedupcausedbytheL1cacheisnotashighasseenin L2 caches, only counting the latency of the network as part of
theFLWbenchmark. the L2 cache latency. However, as we alter the network in our
Impact of increased frequency. One major change in the experiments,asshowninFigure9,webelievethatnetworklatency
RDNAGPUsistheincreasedfrequency.Whiletheincreaseisnot shouldbelowerwhentheL1cacheispresent.Sinceitishardto
anarchitecturalfeature,butinsteadabenefitofshrinkingtransistor usemicrobenchmarkstoseparateoutthelatencyofthenetwork
sizes,weareeagertounderstandhowmuchtheincreasedfrequency versustheL2cacheforbothcases,wemaylosesomeopportunities
canimproveperformance. toidentifythetruebenefitsoftheL1cache.
Here,weuseP3asabaselineandincreasetheGPUcorefre- Limitations. WhilewehaveshownthatNaviSimcanachieve
quencyto1845MHz,whichisthefrequencythattheRX5500XT highaccuracy,thereareseveraluniquearchitecturalfeaturesthat
GPUrunsat.Notethatasweincreasethecorefrequency,wealso arenotpresentlymodeled.Forexample,wedonotsupportthe
increasetheclockspeedthatcontrolstheL0,L1,andL2cachesto CLAUSEinstructionintheRDNAISA[4],whichservesasaperfor-
1845MHz.InAMDGPUs,thecoresandthecachesworkinthesame mancehinttopreventtheDCUsfromswitchingcontexts.Wealso
frequencydomain.WedonotchangetheDRAMconfiguration,so donotimplementinstructionsthatexplicitlyflushtheL0caches.
theDRAMlatencyandbandwidthremainunchanged. Wefindthesefeaturesarerarelyusedanddonothaveamajorper-
Overall,weseetheperformanceofallthebenchmarksimprove, formanceimpact.Wewillbeimplementingbothofthesefeatures
withtherateofchangerangingfrom1.25×to1.84×(seeFigure11). infuturework.
Thissuggeststhattheincreaseinfrequencyleadstoamarkedim- Additionally,NaviSimdeliversasimulatormodelforonlyone
provementinGPUperformance.Thisisparticularlyevidentinthe generationofAMDGPUs.Wedonotsupportsimulatingother
PRandReLUbenchmarks,whichexperienceaspeedupwhichis versionsoftheGPUsandGPUsfromothervendors.However,we
closelycorrelatedwiththeincreaseinthecorefrequency.Thissug- donotconsiderthisasadisadvantageasfocusingononesimulator
geststhattheperformanceofthesetwobenchmarksisdominated modelallowsustoensurethefidelityoftheGPU.Mostwidelyused
bytheGPUfrequency.Afterourfurtherinvestigationwiththe GPUsimulators(e.g.,GPGPUSim[7],Multi2Sim[38])startedwitha
Daisen[36]visualizationtool,wefindthatthistrendisduetoa singleGPUmodelandgraduallyaddednewmodelstothesimulator
fasterwork-groupdispatchingrate.Thesetwobenchmarkshavea infrastructures.Also,consideringthatNaviSimisdevelopedusing
rathershortwork-groupexecutiontime,sothedispatchercannot the same underlying simulator engine as MGPUSim, users can
catchupwiththework-groupretirementspeed.BecausetheCU easilycombinethetwosimulatorstosimulatebothAMDGCN
cancompletethework-groupsfasterthantherateofnewwork- architectureandRDNAarchitecture.
groups arriving to start execution, the number of work-groups
concurrentlyexecutingineachCUislimited.Thisalsolimitsthe
numberofconcurrentin-flightmemorytransactions.Increasing
8 RELATEDWORK
thecorefrequencyacceleratesthedispatchingspeedandincreases
theoccupancyoftheCUs,allowingtheCUstobetterutilizethe GPUSimulators. GPUsimulatorshavebeencriticalinfrastruc-
memorybandwidth.Thiseffectisalsoobservedinthecasestudy turesthatenableGPUarchitecturedesignvalidation.Todate,the
discussedintheDaisenpaper[36].Meanwhile,wealsonoticethat GPUarchitectureresearchcommunityhasdevotedmajoreffort
theaveragespeedup(geometricalmean)isonly1.5×,suggesting intodevelopingGPUsimulatorsandemulators.Earliertoolsin-
thattheincreasedcorefrequencystillneedstobeaccompaniedby cludeBarra[9]andGPUOcelot[12,21]),whichprovidefunctional
improvedDRAMperformance. GPU emulation support, though do not support time modeling.
Next,GPGPUSim[7]andMulti2Sim[38]wereintroducedtode-
liverreliableperformancemodelingofNVIDIAandAMDGPUs.
7 DISCUSSION
GPGPUSimhasbeenextendedtosupportadditionalGPUfeatures,
Usingasimulatortounderstandtherationalebehindarchi- suchasvirtualaddressing[6],concurrentkernelexecution[42],par-
tecturalchanges. OurexperimentsinSection6helpusunder- allelsimulation[16],andtrace-basedsimulation[22].Meanwhile,
standpartofthedesignrationalebehindchangestoAMDGPUs. Multi2SimaddedsupportfortheNVIDIAKeplerarchitecture[14].
Tothebestofourknowledge,ourstudyisthefirsttoperformthis Inrecentyears,moresimulatorshavebeendeveloped,primarily
typeofanalysis.Here,wewanttoreflectonthebenefitsofhaving tosupportnewerGPUarchitectures.Forexample,Macsim[24]sim-
asimulatorsuchasNaviSimavailableforevaluatingthereasons ulatestheIntelGPUarchitecture.TheAMDgem5GPUmodel[15]
forthesechanges. isacomponentaddedtothegem5simulatorandisdedicatedtothe
Overall,wefindthattheflexibilityofusingasimulatorenablesus AMDGCN3architecture.MGPUSim[33]isalsoahigh-performance
toevaluatespecifichardwaredesignchoices.Weareabletoisolate parallel GPU simulator targeting the AMD GCN3 architecture.
343

PACT’22,October8–12,2022,Chicago,IL,USA Bao,Sun,Feric,Shen,Weston,Abellán,Baruah,Kim,Joshi,andKaeli
Moreover,Accel-Sim[22]isanextensiontotheGPGPUSiminfras- [4] AMDInc.2020. "RDNA1.0"InstructionSetArchitecture,ReferenceGuide.
tructurethatcansimulateclosed-sourceGPUprograms,utilizinga https://developer.amd.com/wp-content/resources/RDNA_Shader_ISA.pdf
[5] AMDInc.2022.HIPProgrammingGuide. https://rocmdocs.amd.com/en/latest/
trace-basedsimulationmethod.Accel-Simalsoaddedsupportfor
Programming_Guides/HIP-GUIDE.html
theNVIDIAKepler,Pascal,Volta,andTuringarchitectures.Finally, [6] RachataAusavarungnirun,JoshuaLandgraf,VanceMiller,SaugataGhose,Jayneel
NVArchSim[40]isaninternalsimulatorusedbyNVIDIA,with Gandhi,ChristopherJRossbach,andOnurMutlu.2017.Mosaic:aGPUmemory
managerwithapplication-transparentsupportformultiplepagesizes.InProceed-
muchhigherperformanceascomparedtoGPGPUSim. ingsofthe50thAnnualIEEE/ACMInternationalSymposiumonMicroarchitecture.
ThedesignanddevelopmentofNaviSimhavebeeninspiredby IEEE,Cambridge,MA,USA,136–150.
[7] AliBakhoda,GeorgeLYuan,WilsonWLFung,HenryWong,andTorMAamodt.
existingsimulators.NaviSimsharesthesimulatorcoretechnology
2009.AnalyzingCUDAworkloadsusingadetailedGPUsimulator.In2009IEEE
withMGPUSim,andhence,inheritsthehighperformance,high InternationalSymposiumonPerformanceAnalysisofSystemsandSoftware.IEEE,
flexibility,andmulti-GPUsimulationcapabilityofMGPUSim.Addi- IEEE,Boston,MAUSA,163–174.
[8] Tsann-BimChiou,AlekCChen,MirceaDusa,andShih-EnTseng.2016.Impact
tionally,tothebestofourknowledge,NaviSimisthefirstsimulator
ofEUVpatterningscenarioondifferentdesignstylesandtheirgroundrulesfor
thatcansimulateoneofthenewestRDNAGPUarchitecturesand 7nm/5nmnodeBEOLlayers.InDesign-Process-TechnologyCo-optimizationfor
isvalidatedwithoneofthemostrigorousvalidationprocesses, ManufacturabilityX,Vol.9781.InternationalSocietyforOpticsandPhotonics,
SPIE,Bellingham,WashingtonUSA,978107.
providingareliablebaselineforfutureGPUarchitectureresearch. [9] CarolineCollange,MarcDaumas,DavidDefour,andDavidParello.2010.Barra:
PriorworkwithGPUsimulatorshasexploredandemulating AParallelFunctionalSimulatorforGPGPU.In2010IEEEInternationalSymposium
onModeling,AnalysisandSimulationofComputerandTelecommunicationSystems.
capabilitiesofdifferentaspectsofaGPU.Theneedforarchitectural
IEEE,Miami,Florida,USA,351–360. https://doi.org/10.1109/MASCOTS.2010.43
levelsimulatorscontinuouslygrowsasresearchersemploythem [10] AnthonyDanalis,GabrielMarin,CollinMcCurdy,JeremyS.Meredith,PhilipC.
morewithinthedomainofhighperformancecomputing. Roth,KyleSpafford,VinodTipparaju,andJeffreyS.Vetter.2010.TheScalable
HeterogeneousComputing(SHOC)BenchmarkSuite.InProceedingsofthe3rd
WorkshoponGeneral-PurposeComputationonGraphicsProcessingUnits(GPGPU-
9 CONCLUSION 3).AssociationforComputingMachinery,NewYork,NY,USA,63–74. https:
//doi.org/10.1145/1735688.1735702
Up-to-date and accurate architectural simulators that can faith- [11] ShiDongandDavidKaeli.2017.DNNMark:ADeepNeuralNetworkBenchmark
SuiteforGPUs.InProceedingsoftheGeneralPurposeGPUs(Austin,TX,USA)
fullymodeltoday’scomputingplatformsarekeytoolsetsforboth
(GPGPU-10).AssociationforComputingMachinery,NewYork,NY,USA,63–72.
developing a comprehensive understanding of current comput- https://doi.org/10.1145/3038228.3038239
ingdesigntrends,aswellasevaluatingnewdesignideastobuild [12] NailaFarooqui,AndrewKerr,GregoryDiamos,SudhakarYalamanchili,and
KarstenSchwan.2011.Aframeworkfordynamicallyinstrumentinggpucompute
forward-lookingcomputingplatforms.InthecontextofAMDGPU
applicationswithingpuocelot.InProceedingsoftheFourthWorkshoponGeneral
platforms,thispaperbridgestheseimportantgapsbyproposing PurposeProcessingonGraphicsProcessingUnits.ACM,NewportBeach,CA,1–9.
NaviSim,thefirstcycle-levelsimulatorthattargetsstate-of-the-art [13] DebashisGanguly,ZiyuZhang,JunYang,andRamiMelhem.2020. Adaptive
pagemigrationforirregulardata-intensiveapplicationsunderGPUmemory
RDNA-basedGPUs.Throughourintensiveandrigorousvalidation oversubscription.In2020IEEEInternationalParallelandDistributedProcessing
methodology,whichincludeddevelopingseveralmicrobenchmarks Symposium(IPDPS).IEEE,IEEE,NewOrleans,LouisianaUSA,451–461.
[14] XunGong,RafaelUbal,andDavidKaeli.2017. Multi2SimKepler:Adetailed
andbuildingontheAkitasimulationframework[33],weareable
architecturalGPUsimulator.In2017IEEEInternationalSymposiumonPerformance
tocalibrateNaviSimexecutiontoachieveasmall9.85%average AnalysisofSystemsandSoftware(ISPASS).IEEE,IEEE,SantaRosa,CA,269–278.
simulationerror,ascomparedtohardwareexecutiononanAMD [15] AnthonyGutierrez,BradfordMBeckmann,AlexandruDutu,JosephGross,
MichaelLeBeane,JohnKalamatianos,OnurKayiran,MatthewPoremba,Brandon
RX5500XTRDNAGPU.ToshowcasehowNaviSimcanbeusedto
Potter,SoorajPuthoor,etal.2018.Lostinabstraction:Pitfallsofanalyzinggpus
exploretherationalebehindcriticalarchitecturaldesigndecisions attheintermediatelanguagelevel.In2018IEEEInternationalSymposiumon
made in the transition from pre-RDNA (GCN-based) GPUs, we HighPerformanceComputerArchitecture(HPCA).IEEE,IEEE,Vienna,Austria,
608–619.
carryoutasetofusecasestoquantitativelyanalyzetheeffectsof [16] ClaytonHughes,SimonDavidHammond,MengchiZhang,YechenLiu,Tim
differentRDNAfeaturesonworkloadperformance.Bysupporting Rogers,andRobertJHoekstra.2021.SST-GPU:AScalableSSTGPUComponentfor
PerformanceModelingandProfiling.TechnicalReport.SandiaNationalLab.(SNL-
GPUprogramsdevelopedinbothOpenCLandHIP,NaviSimcon-
NM),Albuquerque,NM(UnitedStates).
tributesavaluablesimulationframeworkforfurtherdesign-space [17] OpenSourceIntiative.1980.TheMITLicense.
exploration,andenablestheresearchanddesignofnext-generation [18] CLJermain,GERowlands,RABuhrman,andDCRalph.2016.GPU-accelerated
micromagneticsimulationsusingcloudcomputing.JournalofMagnetismand
GPUsbasedontheRDNAarchitecture.
MagneticMaterials401(2016),320–322.
[19] JEDECJESD250.2017.Graphicsdoubledatarate6(GDDR6)SGRAMstandard.
JEDECSolidStateTechnologyAssociation.
ACKNOWLEDGMENTS
[20] DavidRKaeli,PerhaadMistry,DanaSchaa,andDongPingZhang.2015.Hetero-
geneouscomputingwithOpenCL2.0.MorganKaufmann,Burlington,MA,USA.
ThisworkissupportedinpartbygiftsfromAMD,William&Mary
[21] AndrewKerr,GregoryDiamos,andSudhakarYalamanchili.2012.Gpuapplication
start-upfunding,andprojectgrantPID2020-112827GB-I00funded development,debugging,andperformancetuningwithgpuocelot. InGPU
byMCIN/AEI/10.13039/501100011033. ComputingGemsJadeEdition.Elsevier,Amsterdam,Netherlands,409–427.
[22] MahmoudKhairy,ZheshengShen,TorMAamodt,andTimothyGRogers.2020.
Accel-Sim:AnextensiblesimulationframeworkforvalidatedGPUmodeling.In
REFERENCES 2020ACM/IEEE47thAnnualInternationalSymposiumonComputerArchitecture
(ISCA).IEEE,IEEE,Valencia,Spain,473–486.
[1] AmirAli Abdolrashidi, Hodjat Asghari Esfeden, Ali Jahanshahi, Kaustubh [23] GwangsunKim,MinseokLee,JiyunJeong,andJohnKim.2014.Multi-GPUSys-
Singh,NaelAbu-Ghazaleh,andDanielWong.2021. Blockmaestro:Enabling temDesignwithMemoryNetworks.InProceedingsofthe47thAnnualIEEE/ACM
programmer-transparenttask-basedexecutioninGPUsystems.In2021ACM/IEEE InternationalSymposiumonMicroarchitecture(Cambridge,UnitedKingdom)
48thAnnualInternationalSymposiumonComputerArchitecture(ISCA).IEEE,IEEE, (MICRO-47).IEEEComputerSociety,USA,484–495. https://doi.org/10.1109/
NewYork,NY,333–346. MICRO.2014.55
[2] AMDInc.2012. AMDGraphicsCoreNextArchitecture. https://www. [24] HyesoonKim,JaekyuLee,NageshBLakshminarayana,JaewoongSim,Jieun
techpowerup.com/gpu-specs/docs/amd-gcn1-architecture.pdf Lim,andTriPho.2012.Macsim:Acpu-gpuheterogeneoussimulationframework
[3] AMDInc.2019.IntroducingRDNAArchitecture,TheallnewRadeongaming userguide.GeorgiaInstituteofTechnology,Atlanta,GA.
architecturepowering“Navi”. https://www.amd.com/system/files/documents/ [25] JoonyoungKimandYounsuKim.2014.HBM:Memorysolutionforbandwidth-
rdna-whitepaper.pdf hungryprocessors.In2014IEEEHotChips26Symposium(HCS).IEEE,IEEE,
344

NaviSim:AHighlyAccurateGPUSimulatorforAMDRDNAGPUs PACT’22,October8–12,2022,Chicago,IL,USA
Cupertino,CA,1–24. 10.1109/IISWC.2016.7581262
[26] AlexKrizhevsky,IlyaSutskever,andGeoffreyE.Hinton.2012.ImageNetClassi- [35] Y.Sun,S.Mukherjee,T.Baruah,S.Dong,J.Gutierrez,P.Mohan,andD.Kaeli.
ficationwithDeepConvolutionalNeuralNetworks.InProceedingsofthe25th 2018.EvaluatingPerformanceTradeoffsontheRadeonOpenComputePlatform.
InternationalConferenceonNeuralInformationProcessingSystems-Volume1 In2018IEEEInternationalSymposiumonPerformanceAnalysisofSystemsand
(LakeTahoe,Nevada)(NIPS’12).CurranAssociatesInc.,RedHook,NY,USA, Software(ISPASS).IEEEComputerSociety,LosAlamitos,CA,USA,209–218.
1097–1105. https://doi.org/10.1109/ISPASS.2018.00034
[27] ShangLi,ZhiyuanYang,DhirajReddy,AnkurSrivastava,andBruceJacob.2020. [36] YifanSun,YixuanZhang,AliMosallaei,MichaelDShah,CodyDunne,and
DRAMsim3:acycle-accurate,thermal-capableDRAMsimulator.IEEEComputer DavidKaeli.2021.Daisen:AFrameworkforVisualizingDetailedGPUExecution.
ArchitectureLetters19,2(2020),106–109. EurographicsConferenceonVisualization40,3(2021),239–250.
[28] MauroMorsianiandRenzoDavoli.1999.Learningoperatingsystemsstructure [37] TheGoProject.2019.EffectiveGo.https://golang.org/doc/effective_go.html.
andimplementationthroughtheMPScomputersystemsimulator.ACMSIGCSE [38] RafaelUbal,ByunghyunJang,PerhaadMistry,DanaSchaa,andDavidKaeli.2012.
Bulletin31,1(1999),63–67. Multi2Sim:ASimulationFrameworkforCPU-GPUComputing.InProceedingsof
[29] JasonJongKyuPark,YongjunPark,andScottMahlke.2017.DynamicResource the21stInternationalConferenceonParallelArchitecturesandCompilationTech-
ManagementforEfficientUtilizationofMultitaskingGPUs.InProceedingsofthe niques(Minneapolis,Minnesota,USA)(PACT’12).AssociationforComputingMa-
Twenty-SecondInternationalConferenceonArchitecturalSupportforProgramming chinery,NewYork,NY,USA,335–344. https://doi.org/10.1145/2370816.2370865
LanguagesandOperatingSystems(Xi’an,China)(ASPLOS’17).Associationfor [39] HenkAVanderVorst.1992. Bi-CGSTAB:Afastandsmoothlyconverging
ComputingMachinery,NewYork,NY,USA,527–540. https://doi.org/10.1145/ variantofBi-CGforthesolutionofnonsymmetriclinearsystems.SIAMJournal
3037697.3037707 onscientificandStatisticalComputing13,2(1992),631–644.
[30] Louis-NoëlPouchetetal.2012.Polybench:Thepolyhedralbenchmarksuite. [40] OresteVilla,DanielLustig,ZiYan,EvgenyBolotin,YaoshengFu,Niladrish
[31] AMDStaff.2014.OpenclandtheAMDAppSDKv2.4. Chatterjee,NanJiang,andDavidNellans.2021. Needforspeed:Experiences
[32] YifanSun,NicolasBohmAgostini,ShiDong,andDavidR.Kaeli.2019.Summa- buildingatrustworthysystem-levelGPUsimulator.In2021IEEEInternational
rizingCPUandGPUDesignTrendswithProductData.CoRRabs/1911.11313 SymposiumonHigh-PerformanceComputerArchitecture(HPCA).IEEE,IEEE,
(2019),1–5. Seoul,Korea(South),868–880.
[33] YifanSun,TrinayanBaruah,SaifulA.Mojumder,ShiDong,XiangGong,Shane [41] VasilyVolkovandJamesWDemmel.2008.BenchmarkingGPUstotunedense
Treadway,YuhuiBao,SpencerHance,CarterMcCardwell,VincentZhao,Harrison linearalgebra.InSC’08:Proceedingsofthe2008ACM/IEEEconferenceonSuper-
Barclay,AmirKavyanZiabari,ZhongliangChen,RafaelUbal,JoséL.Abellán, computing.IEEE,IEEE,Austin,TX,USA,1–11.
JohnKim,AjayJoshi,andDavidKaeli.2019.MGPUSim:EnablingMulti-GPU [42] HaonanWang,FanLuo,MohamedIbrahim,OnurKayiran,andAdwaitJog.2018.
PerformanceModelingandOptimization.InProceedingsofthe46thInternational Efficientandfairmulti-programminginGPUsviaeffectivebandwidthman-
SymposiumonComputerArchitecture(Phoenix,Arizona)(ISCA’19).Association agement.In2018IEEEInternationalSymposiumonHighPerformanceComputer
forComputingMachinery,NewYork,NY,USA,197–209. https://doi.org/10. Architecture(HPCA).IEEE,IEEE,Vienna,Austria,247–258.
1145/3307650.3322230 [43] LingyuanWang,MiaoqingHuang,andTarekEl-Ghazawi.2011. Exploiting
[34] Y.Sun,X.Gong,A.Ziabari,L.Yu,X.Li,S.Mukherjee,C.Mccardwell,A.Villegas, concurrentkernelexecutionongraphicprocessingunits.In2011International
andD.Kaeli.2016.Hetero-mark,abenchmarksuiteforCPU-GPUcollaborative ConferenceonHighPerformanceComputing&Simulation.IEEE,IEEE,Istanbul,
computing.In2016IEEEInternationalSymposiumonWorkloadCharacterization Turkey,24–32.
(IISWC).IEEEComputerSociety,LosAlamitos,CA,USA,1–10. https://doi.org/
345
