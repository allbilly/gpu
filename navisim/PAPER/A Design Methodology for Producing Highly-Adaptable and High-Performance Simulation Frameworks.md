A Design Methodology for Producing Highly-Adaptable and
High-Performance Simulation Frameworks
ADissertationPresented
by
YuhuiBao
to
TheDepartmentofElectricalandComputerEngineering
inpartialfulfillmentoftherequirements
forthedegreeof
DoctorofPhilosophy
in
ElectricalandComputerEngineering
NortheasternUniversity
Boston,Massachusetts
December2024

Contents
ListofFigures iv
ListofTables vi
AbstractoftheDissertation vii
1 Introduction 1
1.1 WhyDoestheCommunityContinuallyNeedNewandAdaptableSimulators . . . 2
1.1.1 DynamicsintheGPUMarket . . . . . . . . . . . . . . . . . . . . . . . . 2
1.1.2 AcceleratingChallengingWorkloads . . . . . . . . . . . . . . . . . . . . 3
1.1.3 CommonIssueswithSimulators . . . . . . . . . . . . . . . . . . . . . . . 5
1.2 ChallengesinDesigningaGPUSimulator: ImplicationsofPreviousConsiderations 6
1.3 TheSpeedofSimulation . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 7
1.4 ContributionsofThisThesis . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 9
1.5 TheOrganizationofthisProposal . . . . . . . . . . . . . . . . . . . . . . . . . . 10
2 Background 11
2.1 GPUProgrammingFrameworks . . . . . . . . . . . . . . . . . . . . . . . . . . . 11
2.2 InstructionSetArchitectureversusArchitectureversusMicroarchitecture . . . . . 12
2.3 GPUMicroarchitecture . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 13
2.3.1 KernelDispatching . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 13
2.3.2 ComputeUnit . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 14
2.3.3 TheGPUMemoryHierarchy . . . . . . . . . . . . . . . . . . . . . . . . . 15
2.4 ArchitecturalSimulators . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 16
2.4.1 MicroarchitectureSimulation . . . . . . . . . . . . . . . . . . . . . . . . 16
2.4.2 ClassificationofSimulators . . . . . . . . . . . . . . . . . . . . . . . . . 17
2.5 PerformanceEvaluationTools . . . . . . . . . . . . . . . . . . . . . . . . . . . . 18
2.5.1 perf . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 19
2.5.2 PAPI . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 19
2.5.3 pprof . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 19
2.6 FullyHomomorphicEncryption . . . . . . . . . . . . . . . . . . . . . . . . . . . 19
2.6.1 TheCKKSFHEScheme . . . . . . . . . . . . . . . . . . . . . . . . . . . 20
2.6.2 BuildingBlocksoftheCKKSFHEScheme . . . . . . . . . . . . . . . . . 22
i

3 RelatedWork 27
3.1 GPUSimulators . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 27
3.1.1 Commercially-developSimulators . . . . . . . . . . . . . . . . . . . . . . 27
3.1.2 AcademicResearchSimulators . . . . . . . . . . . . . . . . . . . . . . . 28
3.2 ExtendingGPUArchitecturestoAccelerateFHE . . . . . . . . . . . . . . . . . . 29
3.2.1 Extending/ModifyingCurrentGPUArchitectures . . . . . . . . . . . . . . 29
3.2.2 AcceleratingFHE . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 30
4 NaviSim 32
4.1 DesignGoals . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 33
4.2 SimulatorDevelopmentMethodology . . . . . . . . . . . . . . . . . . . . . . . . 34
4.3 RDNAArchitecture . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 36
4.4 ImplementationofNaviSim . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 38
4.4.1 ISAEmulation . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 38
4.4.2 WavefrontDispatching . . . . . . . . . . . . . . . . . . . . . . . . . . . . 38
4.4.3 DCUmodeling . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 39
4.4.4 MemoryHierarchy . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 40
4.4.5 UserInterface . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 40
4.5 Evaluation . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 41
4.5.1 SimulatorValidation . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 41
4.6 CaseStudies: UnderstandingtheRDNAArchitectureFeatures . . . . . . . . . . . 46
4.6.1 ImpactofchangesintheISAandtheintroductionoftheDCU . . . . . . . 47
4.6.2 ImpactoftheL1cache . . . . . . . . . . . . . . . . . . . . . . . . . . . . 49
4.6.3 Impactofincreasedfrequency . . . . . . . . . . . . . . . . . . . . . . . . 50
5 Design Space Exploration: Accelerating FHE with Microarchitectural Extensions
(GME) 52
5.1 ChallengesofAcceleratingCKKSFHE . . . . . . . . . . . . . . . . . . . . . . . 52
5.2 ExploringMicroarchitecturalExtensionswithNaviSim . . . . . . . . . . . . . . . 54
5.2.1 cNoC:CU-sideinterconnect . . . . . . . . . . . . . . . . . . . . . . . . . 55
5.2.2 EnhancingtheVectorALU . . . . . . . . . . . . . . . . . . . . . . . . . . 57
5.2.3 LABS:Locality-AwareBlockScheduler . . . . . . . . . . . . . . . . . . . 59
5.3 Evaluation . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 61
5.3.1 Methodology . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 61
5.3.2 Results . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 63
6 SimulationPerformanceEvaluationandOptimization 68
6.1 Methodology . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 69
6.1.1 Generalsimulatordesign . . . . . . . . . . . . . . . . . . . . . . . . . . . 69
6.1.2 Parallelsimulationscalability . . . . . . . . . . . . . . . . . . . . . . . . 71
6.1.3 PerformanceMetrics . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 72
6.2 Evaluation . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 72
6.2.1 ExperimentalSetupandProfilingMethods . . . . . . . . . . . . . . . . . 72
6.2.2 PerformanceAnalysis . . . . . . . . . . . . . . . . . . . . . . . . . . . . 75
6.3 PerformanceImprovementResults . . . . . . . . . . . . . . . . . . . . . . . . . . 83
ii

6.4 Summary . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 90
7 ConclusionandFutureWork 93
7.1 Conclusion . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 93
7.2 FutureWork . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 94
Bibliography 96
A FullSynchronizationBlockingProfiles 107
iii

List of Figures
1.1 The roofline model for the FHE operation with size of the ciphertext N =216 on
AMDMI100GPU . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 4
1.2 Acomparisonofasimulator’sreportedkernelexecutiontimeandthenativeRX5500
XTGPUexecutiontimewhenrunningtheBICGbenchmarkfromPolyBench[1]. 5
1.3 SimulationslowdownofbenchmarksrunningonNaviSimunderemulationmode
andtimingsimulationmode . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 8
2.1 TheMicrorchitectureofCDNAGPUs . . . . . . . . . . . . . . . . . . . . . . . . 14
2.2 TheMicrorchitectureofCDNAComputeUnit . . . . . . . . . . . . . . . . . . . . 15
2.3 FHE offers a safeguard against online eavesdroppers as well as untrusted cloud
servicesbyallowingdirectcomputationonencrypteddata. . . . . . . . . . . . . . 20
2.4 TheoperationflowoftheCKKSscheme. . . . . . . . . . . . . . . . . . . . . . . 21
2.5 PolynomialMultiplication. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 23
4.1 AComparisonbetweentheGCNArchitecture[2]andtheRDNAArchitecture[3]. 37
4.2 TheinstructionpipelinemodelinaNaviSimDCU.EachDCUhasfourcopiesofthe
structureshown. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 39
4.3 ComparisonofthePointerChasingmicrobenchmarkbetweenNaviSimandtheRX
5500XThardware. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 45
4.4 Simulator validation against the AMD RX 5500 XT GPU. The x-axis plots the
problem size and the two y-axes plot the kernel execution time and relative error.
ThenumbersintheparenthesesrepresenttheaveragedifferencebetweenNaviSim
simulationandthehardwareexecutionforeachapplication. . . . . . . . . . . . . . 46
4.5 SimulatorvalidationagainsttheAMDRadeonProW6800GPU.Thex-axisplots
theproblemsizeandthetwoy-axesplotthekernelexecutiontimeandrelativeerror.
ThenumbersintheparenthesesrepresenttheaveragedifferencebetweenNaviSim
simulationandthehardwareexecutionforeachapplication. . . . . . . . . . . . . . 47
4.6 TheeffectofchangingfromCUstoDCUsandtheISAfromGFX803(GCN3)to
GFX1010(RDNA)onexecutiontime. . . . . . . . . . . . . . . . . . . . . . . . . 48
4.7 Theeffectofaddingthenewwrite-evictL1cacheonbenchmarkexecutiontime. . 49
4.8 Thecachehitrateofeachlevelofcache,beforeandaftertheL1cacheisadded. . . 50
4.9 Theeffectoftheincreasedcorefrequencyontheoverallperformance. Asweincrease
thecorefrequencyby1.845×,theperformanceimprovedby≈1.5×. . . . . . . . . 51
iv

5.1 NTTkernelsprofileofALUutilization,cachehitratesandmemoryintensity . . . 53
5.2 Inter-CUcommunication: Traditional(left)versusproposedon-chipnetwork. . . . 55
5.3 Proposedhierarchicalon-chipnetworkfeaturingaconcentrated2Dtorustopology. 56
5.4 Influenceofindividualproposedmicro-architecturalextensiononarchitecturalper-
formancemetrics. Metricsillustrateacumulativeprofilewhereeachenhancement
buildsupontheprecedingsetofimprovements . . . . . . . . . . . . . . . . . . . 64
5.5 Speedup achieved from each micro-architectural extension. The baseline refers
to a vanilla MI100 GPU. The reported speedup is cumulative, with each micro-
architecturalenhancementbuildinguponthepreviousones . . . . . . . . . . . . . 65
5.6 Exploringtheimpactofon-chipmemorysizeonFHEworkloadperformance. . . . 67
6.1 TheCPUprofilefromtheGopproftool,whilerunningtheFIRbenchmarktiming
simulation. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 76
6.2 Thespeedupachievedthroughparallelsimulationwhenrunningemulationmode. . 77
6.3 Executiontimebreakdowningoroutineanalysisview,usingtheGotracetool. . . . 78
6.4 ThesynchronizationblockingprofilefromtheGotracetool,whilerunningtheReLU
benchmarkparallelemulation. . . . . . . . . . . . . . . . . . . . . . . . . . . . . 79
6.5 Thespeedupachievedthroughparallelsimulationwhenrunningtimingsimulation. 80
6.6 ThesynchronizationblockingprofilefromGotracetoolforReLUbenchmarkparallel
timingsimulation . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 81
6.7 Executionperformancecomparisoncomparingthebaselineandimprovementsmade
totheparallelengine. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 82
6.8 Thesynchronizationblockingprofileafterapplyingthereversing-orderstrategy. . . 84
6.9 Performancecomparisonbetweenthebaseline,themodifiedComputeUnitandthe
newWorkgroupDispatcher.. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 86
6.10 EmulationperformancecomparisonbetweentheoriginalALUandthenewALU
implementations. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 87
6.11 Timing performance comparison between the original ALU and the new ALU
implementations. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 88
6.12 Thespeedupachievedthroughparallelsimulationwhenrunningemulationmodeon
platform1. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 89
6.13 Thespeedupachievedthroughparallelsimulationwhenrunningemulationmodeon
platform2. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 90
A.1 ThesynchronizationblockingprofilefromGotracetoolforReLUbenchmarkparallel
emulation. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 108
A.2 ThesynchronizationblockingprofilefromGotracetoolforReLUbenchmarkparallel
timingsimulation. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 109
A.3 Thesynchronizationblockingprofileafterapplyingthereversing-orderstrategy. . . 110
v

List of Tables
2.1 TheHEbuildingblocksusingCKKS. . . . . . . . . . . . . . . . . . . . . . . . . 24
2.2 TheCKKSparametersanddescriptions. . . . . . . . . . . . . . . . . . . . . . . . 26
2.3 PracticalparametersforourFHEoperations.. . . . . . . . . . . . . . . . . . . . . 26
4.1 ThelistofmicrobenchmarksusedtocalibrateNaviSim. . . . . . . . . . . . . . . . 35
4.2 Theconfigurationsofvalidationplatforms. . . . . . . . . . . . . . . . . . . . . . . 41
4.3 Simulatorconfiguration . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 42
4.4 FullBenchmarks . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 44
5.1 Cyclecountsfor64-bitmodulusinstructionscomparingMODandWMACfeatures 59
5.2 MI100GPUParameters . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 62
5.3 PerformanceofvariousFHEbuildingblocks . . . . . . . . . . . . . . . . . . . . . 63
6.1 Benchmarks . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 73
6.2 Simulatorconfiguration . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 74
6.3 Evaluationplatform . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 74
6.4 Stepstotunetheperformanceofasimulator . . . . . . . . . . . . . . . . . . . . . 91
vi

Abstract of the Dissertation
A Design Methodology for Producing Highly-Adaptable and
High-Performance Simulation Frameworks
by
YuhuiBao
DoctorofPhilosophyinElectricalandComputerEngineering
NortheasternUniversity,December2024
Dr. DavidKaeli,Advisor
Computerarchitecturesimulatorsplayanessentialroleinthedevelopmentandoptimization
ofcomputerhardware. Avarietyofsimulatorshavebeendevelopedtoexplorethedesignspaceof
CPUs,GPUs,andcustomeraccelerators. AsGPUscontinuetogrowinpopularityforaccelerating
demandingapplications,suchashigh-performancecomputingandmachinelearning,GPUarchitects
havebeenpushingtheenvelopeofGPUperformanceineverynewGPUgeneration. GPUvendors
(e.g.,NVIDIAandAMD)havebeenintroducingsubsequentgenerationsofGPUarchitecturesand
productswithupdatedinstructionsetarchitectures(ISAs)andnewmicroarchitecturalfeaturesevery
2-3years. Modelingthestate-of-the-artarchitectureisacrucialfeatureofGPUsimulators,whichare
usedtocharacterizeandacceleratechallengingworkloadsfacilitatingperformanceevaluationand
designexploration. However,theeffortrequiredtodesignandconstructanaccurateandperformant
simulator is huge. Due to the rapid rate of innovation in GPU technology, any simulator that is
over-customizedtocapturethedesignofaspecificarchitecturewillquicklybecomeoutdated. Thus,
weneedtodevelopadesignmethodologyforsimulatorsthatcanguardagainstthistrend,embracing
futurearchitectures.
Inthisdissertation,weproposeadesignmethodologyforproducinghighly-adaptableand
high-performancesimulationframeworks. Weaimtodesignsimulatorsfeaturinghighadaptability,
beingabletoaccommodatefuturealterationsorextensions,highperformanceandhighfidelity. We
leveragetheAkitasimulatorframeworktoenablethemodularandextensibledesignofvariousGPU
components. To fulfill the goal of high fidelity, we design a set of microbenchmarks to evaluate
individualGPUsubsystems. Wedemonstratehowwefollowourdesignmethodologytoachieve
a highly-adaptable and accurate simulator — NaviSim, which provides the flexibility to support
simulation of three different ISAs. To demonstrate the full utility of the NaviSim simulator, we
vii

conduct a performance study of the impact of individual architecture features revealing the high
flexibilityandconfigurabilityofNaviSim. Inaddition,weshowcasehowNaviSim’shighadaptability
contributestodesignspaceexploration,offeringsolutionstoenhancetheperformanceofreal-world
demandingapplications.
Fast simulation speed is one of the key requirements of any simulators. NaviSim is
designed to support multi-threaded execution, which is able to leverage the parallel capabilities
offered by today’s multi-core CPUs, enabling parallel simulation. In this thesis we identify key
performancebottlenecksintermsofbothserialandparallelsimulationexecutemodesandoptimize
simulation speed. We also present lessons learned about efficient simulator design and provide
guidanceforfuturesimulatordevelopers.
viii

Chapter 1
Introduction
Overthepast20years,GPU(GraphicsProcessingUnits)hasbecomeanessentialpower
that drives the rapid growth of many scientific and engineering applications, revolutionizing the
landscapeofscientificandengineeringadvances. Comparedwiththetraditionalmulti-coreCentral
ProcessingUnit(CPU),GPUsareconsideredtobethroughput-optimizeddevicewiththousandsof
computecoresandhighmemorybandwidth. Applicationsthatexhibitahighdegreeofdata-level
parallelismcanbenefitgreatlyfromtheimpressiveperformancescalingdeliveredbyGPUs. Thisis
especiallyintheartificialintelligence(AI)[4],whereGPUshavebeenusedtodrasticallyreduce
trainingtimesforDNNs(DeepNeuralNetworks)[5]. GPUsarenotonlyaneffectiveaccelerator
for AI workload, many data-driven applications have been accelerated using a GPU. As a result,
GPUshavebeenwidelydeployedaspartof168supercomputersfoundintheTop500supercomputer
list[6]. GPUshavealsocontributedtothewidespreaddeploymentofnumerousedgecomputing
platforms,embeddedsystems,andportablepersonaldevices.
TracingbackthehistoryofGPU,thetermGPGPUs(GeneralPurposeGPUs)markeda
significantmilestoneineffortstobroadlyperformgeneral-purposecomputationleveringtheparallel
processingpowerofGPUs. GPUswereoriginallydesignedtorender3Dgraphics,thoughtheyhave
beenre-architectedtosupportgeneral-purposecomputationashigh-performanceacceleratorsfor
a growing range of markets. This widespread adoption of GPUs as the accelerator of choice has
beenpartlydrivenbythedevelopmentofmoreuser-friendlyandefficientprogrammingframeworks
suchasCUDA[7],HIP[8,9]andOpenCL[10]. Thesehigh-levelprogramminglanguagesprovidea
well-definedsyntax,aswellasoptimizingcompilers,drivers,andruntimestoallowdevelopersto
effectivelyexploittheparallelarchitectureofGPUs.
AsGPUsevolvedbeyondservingonlyasgraphicsdevice,theyhavebeenextensivelyused
1

CHAPTER1. INTRODUCTION
for various types of compute-intensive workloads that can benefit from their parallel processing
power. GPU-acceleratedlinearalgebralibraries,suchascuBLASinCUDAforNVIDIAGPUs,and
hipBLASfromAMD’sROCmlibraries,arewidelyusedincompute-intensiveapplicationssupporting
accelerationacrossarangeoffields. Inaddition,scientificsimulationsoftenutilizeGPUsforcomplex
computations. Manyscientificsimulationworkloads,suchasweatherforecasting,fluiddynamics,
and molecular modeling, involve the execution of a large number of mathematical operations,
workingwithextensivedatasets,runningcomplexalgorithmsthatconsumeaconsiderableamount
ofprocessingtimeandpower. GPUscanacceleratethesecomputationsbyexecutingthousandsof
computationssimultaneously,whilestreamingdataconcurrently,thusreducingtheoverallprocessing
time.
Computerarchitecturesimulatorsplayanessentialpartinthedevelopmentandoptimization
ofcomputerhardware[11]. Duetothehugecostassociatedwithbuildingrealhardwarefortesting
andverification,computerarchitectureresearchersoftenrelyonsimulationandmodelingtechniques
to evaluate different design options. Architectural simulators are commonly used to guide early
designspaceexploration,performanceoptimization,andpre-siliconverification,supportingmore
efficient and cost-effective development cycles. Having access to an architectural simulator that
delivershighfidelity,highflexibility,andhighperformanceisessentialforGPUarchitectureresearch.
However, building a accurate and performant simulator is a very challenging task. Designing an
architecturalsimulatorthatcanaccuratelyreproducethebehaviorofthehardwarebeingmodeled,is
easytoextendormodify,andcansatisfysimulationperformancerequirementsisamajoreffort.
1.1 Why Does the Community Continually Need New and Adaptable
Simulators
1.1.1 DynamicsintheGPUMarket
Every2-3years,GPUvendors(e.g.,NVIDIAandAMD)havebeenintroducingsubsequent
generationsofGPUarchitecturesandproducts,eachequippedwithavarietyofinnovativedesigns.
While some GPU generations maintain the nearly the same instruction set architecture, others
introduce major architectural overhauls, changing many aspects of the GPU’s instruction set and
systemorganization. TakeforexampleAMD’sGraphicsCoreNext(GCN)architecture[2],whichhas
beenadominantGPUarchitectureforthelastdecade(2011-2019). GCNcontinuedtoevolvethrough
fivegenerations,witheachgenerationintroducingafewnewfeaturesineachgeneration. In2019,
2

CHAPTER1. INTRODUCTION
AMDmovedawayfromtheGCNarchitectureandstartedtorolloutanewRDNAarchitecture[3].
Oneyearlater,AMDannouncedtheirCDNAarchitecture[12],equippedwithaMatrixCoreEngine
tosupportAIandmachinelearningapplications.
AMD’sRDNAarchitecturewasamajorredesign,asAMDmodifiednearlyeveryaspect
of the GPU’s architecture. This included changes to the instruction set architecture (ISA), core
architecture, and memory hierarchy. In the cores, the RDNA architecture reduced the number
of co-scheduled threads (i.e., the wavefront) from 64 to 32 to cope with a higher level of thread
divergenceinmodernworkloads. Thememoryhierarchywasalsoextended,addinganextralayer
of cache between the original L1 and L2 caches. This reduced the burden of the L2 caches and
simplified the massive L1 to L2 network. However, not all exact details of industrial designs are
disclosedandsomehiddenenhancementswerealsomadeunderthehood. Incomparison,theCDNA
architecture mostlyinherits the design ofGCN architecture, but with a brand-new instruction set
architectureandaugmentedprocessingunitstargetingspecializedworkloads.
1.1.2 AcceleratingChallengingWorkloads
AlthoughGPUsprovidegreatpotentialtoaccelerateavarietyofworkloads,somework-
loadsfaceperformanceissueswheretheycannotefficientlyutilizetheparallelcomputingresources
providedontheGPU.Theseperformanceissuescanbeattributedtofactorssuchasthreaddiver-
gence,irregularmemoryaccesspatterns,andotherworkloadcharacteristics,impactingourabilityto
optimizetheperformanceoftheseworkloads.
Any workload that includes inherent parallelism across vectors or threads, and can be
executedconcurrently,canleveragethecapabilitiesofGPUs. However,theamountofeffortrequired
tomodifyanalgorithmtoexposethisparallelismtargetedforGPUexecutioncanvarygreatly. One
classofworkloadthatcanmapsnaturallytoaGPUinvolvesembarrassinglyparallelcode. However,
navigatingreal-worldworkloadsprovestobeachallengingtask,asachievingtheanticipatedideal
performancegainisnotstraightforward. Thiscomplexityisexemplifiedwhenconsideringintricate
workloads such as graph neural networks (GNNs) and homomorphic encryption (HE). In these
instances,theinherentcharacteristicsoftheworkloadspresentchallengesduringmappingtoaGPU.
LetusfocusonHEworkloadsmoreclosely.
Fullyhomomorphicencryption(FHE)offersidealsecuritybyallowingcomputationsto
be run directly on encrypted operands. FHE adopters can offload their encrypted private data to
third-partycloudserviceproviderswhilepreservingend-to-endprivacy. However,onechallengein
3

CHAPTER1. INTRODUCTION
deployingFHEinreal-worldsystemsisthehighcomputationalcostsassociatedwithFHEoperations.
Computationsondataencryptedviastate-of-the-artHEschemes,suchasCKKS[13]andTFHE[14],
experienceaslowdownof4–6ordersofmagnitude,ascomparedtorunningthesamecomputation
onunencrypteddata[15].
ManypriorstudieshavetriedtoaccelerateFHEkernelsbydevelopingGPU-basedimple-
mentationsandapplyingarangeofoptimizations,includingmicroarchitecturalimprovements[16],
low-levelGPUkerneloptimizationtechniques[17],andmulti-GPUoptimization[18]. Eventhough
theseimplementationsoutperformCPU-basedsolutions,theperformancebenefitsoveranout-of-
the-boxGPUimplementationisstilllimited. Weshowarooflinemodel[19]foroneinstanceofthe
FHEwhenranonanAMDMI100GPU.Wecanseethatthereisconsiderableroomforperformance
improvement (See Figure 1.1). Based on the specification of the AMD MI100 GPU, we plot its
8.8
2.2
0.55
0.14
8 16 32 64 128
Arithmetic Intensity (Flops/Byte)
)s/spolFT(
ecnamrofreP
Theoretic Peak:23.1TFlops
MI100 Roofline
baseline
software-level
optimization
Figure1.1: TherooflinemodelfortheFHEoperationwithsizeoftheciphertextN =216 onAMD
MI100GPU
rooflinemodel. Intherooflinemodelplot,weindicatetheperformanceintflops/sandthearithmetic
intensity(AI)inflops/bytemeasurementsofthebaselineandoptimizedGPUimplementationofone
FHEoperation.
DespitethepotentialbenefitsofparallelprocessingonaGPU,FHEoperations,whichisa
4

CHAPTER1. INTRODUCTION
memory-boundworkload,inherentlyfaceslimitationstoexploitthemassivecomputationresources
andhighmemorybandwidth. WithcloserexaminationoftheFHEworkload,wecancharacterizethe
inherentbottlenecks(e.g.,large(dozensofMB)ciphertexts,amassivevolumeofintegerarithmetic
operations,andirregularmemoryaccesspatterns. Theseissuepresentachallengewhenmappingto
existingGPUarchitectures. Historically,GPUshavebeendesignedtoexecutethousandsofthreads
inparallel,exploitingregularmemoryaccesspatternsandrichfloating-pointcomputations.
1.1.3 CommonIssueswithSimulators
Inability to simulate new architectures Most simulators are designed to simulate a
specificinstructionsetarchitecture,typicallysupportingonlyasinglemicroarchitecture.
This lack of up-to-date simulation tools within academia is potentially harmful to the
computerarchitectureresearchcommunity. Researchpapersoftendescribetheirresearchmethod-
ologyusingstatingthat“wemodifiedsimulatorS,whichwasoriginallyvalidatedagainstproduct
A,tomodelproductB,changingthenumberthenumberofcoresandthecachesizes.” Typically,
thisproductBisafewgenerationslaterthanproductA,andthestudyimplicitlyassumesthatthe
simulatorcanstillcorrectlymodelthenewarchitecture. However,thenuancesinthearchitecture
designandparameterselectionmaysignificantlyimpacttheperformancecharacteristics.
600
400
200
0
64 256 512 768 1024
Problem Size
)su(emiT
Simulation Native Execution
Turning point in simulation
but not in native execution
Figure1.2: Acomparisonofasimulator’sreportedkernelexecutiontimeandthenativeRX5500
XTGPUexecutiontimewhenrunningtheBICGbenchmarkfromPolyBench[1].
Lookingatanexample,simplymodifyingthepubliclyavailableparameters(corecount,
frequency,cachesizes,DRAMbandwidth,etc.) ofMGPUSim[20],whichwasvalidatedagainst
a 2015-released R9 Nano GPU to model the RX 5500 XT GPU (RDNA architecture), as shown
inFigure1.2,canresultinsimulationerrorsaslargeas318%whenrunningtheBICGbenchmark
5

CHAPTER1. INTRODUCTION
fromPolyBench[1]. Notonlyaretheperformancediscrepancieshigh,buttheperformancedoes
notfollowatrendswehavemeasuredonRX5500XThardware. Themajorinaccuraciesreported
abovesuggestthatsimplyadjustingpubliclyknownparametersisinsufficienttomodelabrand-new
architecture and can lead to misguided conclusions. A careful redesign and re-calibration of the
simulatorinfrastructureisnecessarytoprovideatrustworthybaselinemodelforthenextgeneration
ofGPUarchitectureresearch.
LackofModularityWhendesigningasimulator,itisimportanttoutilizehighmodularity
ofthecode,asthiswillallowuserstomoreeasilymodifyandextendthesimulatortofittheneedsof
theirintendedstudy. Withoutmodularity,implementingextensionstoasimulatorcanbechallenging
anderrorprone. Themainproblemwillbethelackofindependencebetweensimulatorcomponents.
Users will find it nearly impossible to replace an existing simulator component with their own
component,withoutcausinganyside-effectsinotherpartsofthesimulator.
Thislackofmodularitywillimpactauser’sabilitytoquicklyandeffectivelymodifyor
extendtheGPUsimulator. Additionally,hardwarearchitectswillhavedifficultyperformingdesign
spaceexplorationwhenintroducingnewfeaturestothecurrentdesign. Theywillneedtofrequently
modifythesimulatortoimplementandverifytheirideas.
1.2 Challenges in Designing a GPU Simulator: Implications of Previ-
ous Considerations
In light of the previous discussion, we need to design and formulate a GPU simulator
characterizedbyahighdegreeofadaptability. Thissimulatorshoulddemonstratethecapabilityto
supportmultipleexistingGPUinstructionsetarchitecturesandmicroarchitectures. Thetruevalueof
thesimulatorwilllieinitsabilitytoaccommodatefuturealterationsorextensions.
Weneedtodevelopastrategytoguidesimulatordesignanddevelopment,withthegoal
ofproducingtoolsthatpossessahighdegreeofadaptability. Wefocusonthehighmodularityof
thecodeandtheabilitytoextendthesimulatorwithoutmodifyingit. Togroundourdiscussionin
reality,wehavedevelopedNaviSim[21],anovelGPUsimulatorthatmodelstheAMDRDNAand
CDNAarchitectures. ReflectingonthedesignanddevelopmentprocessweusedforNaviSim,we
highlightthesimulator’sflexibilitytosupportsimulationsofthreedifferentISAsandtwodifferent
microarchitecturesofthreeGPUarchitectrues. NaviSimalsofeatureshighextensibility,allowing
developerstoadd,notonlynewcomponents,includingnewfunctionunitsandinterconnects,and
6

CHAPTER1. INTRODUCTION
newtypesofinstructions,butalsonewISAemulatorsandnewmicroarchitecturalfeatures,without
requiringmajormodifications. Inaddition,weshowcasehowNaviSim’shighadaptabilitycontributes
todesignspaceexploration,offeringsolutionstoenhancetheperformanceofdemandingapplications.
NaviSimcanbeusedtoidentifyhardwarebottlenecksduringworkloadexecutionbyidentifyingthe
microarchitecture-levelbottlenecksandperformingdesignspaceexplorationondifferentarchitectural
featuresfulfillingthespecialneedsoftheintricateworkloads.
1.3 The Speed of Simulation
Over the past decade, the rapid evolution of GPU architectures has been driven by the
demands posed by high-performance computing (HPC) and machine learning (ML) workloads.
Giventhenatureoftheseworkloads(longrunningandpotentiallyhighlymemoryintensive),they
pose challenges when designing a GPU simulator. Ideally, we would like to be able to quickly
explore a wide design space. Rapidly evaluating the impact of a variety of new features across a
diversesetofworkloadsisnecessaryforGPUarchitectstoquantifythepotentialoftheirarchitectural
innovationsforfutureGPUs. Awidelyacknowledgedproblemencounteredwitharchitectural-level
or cycle-level simulators is their slow simulation speed [22]. Usually, simulating even a single
cyclerequiresmuchmoreresourcesandtimethanexecutingonecycleonthehardware,causinga
significantslowdowncomparedwithactualhardware. Thus,simulatingasinglesecondofexecution
couldpotentiallytakedays,weeks,orevenyears. Wequantifyandshowtheslowdownofrunning
differentbenchmarksinNaviSiminFigureFigure1.3. Slowdown,asametric,isdefinedasthewall
clocktimeofthesimulationoverthesimulatedkernelexecutiontime. Thisresultsuggeststhatthe
simulationslowdownfactorvariesamongdifferentbenchmarksandisimpactedbythedesignand
implementationofthesimulation. Onaverage,emulationexhibitsaslowdownofroughly5ordersof
magnitudeandatimingsimulationexperiencesnearlyaslowdownof6ordersofmagnitude.
Ideally,asimulatorshouldatleastsatisfythreedesirableproperties: itshouldbeaccurate,
fast, and execute a wide range of workloads [11]. While fast is not as clearly defined as other
properties, due to the lack of a set of performance metrics to evaluate the architectural simulator
intermsofexecutionspeed. Priorstudiestargetingthespeedissuefocusedonemployingarange
ofapproaches,includingparallelordistributedsimulation[23,24],sampling-basedtechniques[25,
26, 27], and trading off simulator component fidelity for speed [28]. However, these approaches
oftenfailtoconsidertheopportunitieshiddeninthesimulatordesign,overlookingpotentialsoftware
optimizationsthatcouldgreatlyenhancesimulationspeed.
7

CHAPTER1. INTRODUCTION
1.5
1.0
0.5
ATAX BS FWT FIR FLW KM MM MT PR ReLU SPMV
Benchmark
nwodwolS
Emulation Timing Simulation
1e6
Figure1.3: SimulationslowdownofbenchmarksrunningonNaviSimunderemulationmodeand
timingsimulationmode
Inthisthesis,wewilldevelopamethodologyforevaluatingtheperformanceofarchitecture
simulators. Equippedwithourmethodology,weapplyourmethodstooptimizethesimulationspeed
NaviSim. To achieve this goal, we define a set of performance metrics that help evaluate the
simulation speed in terms of various simulated metrics (e.g., kernel execution time, number of
simulatedinstructions,numberofmemorytransactions). Webelieveanextensiveevaluationofthe
simulatorspeedacrossrepresentativeconfigurationscanprovideusefulguidancetobetterunderstand
any bottlenecks present in the simulator, as well as how to remedy them. Next, we conduct a
software-basedprofiling,understandtheimpactofvariousworkloadsonsimulatorthroughput. The
profilingresultshelpustoexaminebettermodelingstrategiesandsoftwareengineeringtechniques
thatwecanemploytoimprovesimulatorthroughput. Basedonourprofilingresults, weperform
performancetuningoftheNaviSimsimulatorfromasoftwareperspective. Sincetoday’ssimulators
are run on multi-core CPU system, we study the parallel simulation mode of NaviSim (NaviSim
isoneofthefewsimulatorsthatismulti-threaded)andimprovetheparallelsimulationscalability.
Wesummarizelessonslearnedaboutefficientsimulatordesignandattempttoprovideguidancefor
futuresimulatordevelopers.
8

CHAPTER1. INTRODUCTION
1.4 Contributions of This Thesis
Thegoalofthisthesisistoproposeadesignmethodologyforproducinghighly-adaptable
andhigh-performancesimulationframeworks. Wedevelopastate-of-the-artGPUsimulatortodemon-
stratehowwefollowthismethodologytodesignanddevelopasimulatorachievinghigh-adaptability
and high performance. We also develop metrics and methods for evaluating the performance of
architecturalsimulators. Thecontributionsofthisthesisinclude:
• We present a simulator development methodology that can be generalized to design any
architecturalsimulatorthatalignswiththegoalsofhighmodularity,extensibility,andfidelity.
Toprovideaconcreteexampleofthegeneraldesignmethodology,wearchitectanewGPU
simulatornamedNaviSim[21],whichcansimulatethreeinstructionsetsandtwodifferent
microarchitectureimplementationsofthreeGPUarchitectures. Giventheinherentflexibibility
in our emulator design, NaviSim enables the modeler to mix and match instruction sets
with different microarchitecture implementations. On top of this feature, we conduct a
performancestudyoftheimpactofindividualarchitecturefeatures,providingahighlyflexible
andconfigurableframework.
• WedemonstratetheadaptabilityofNaviSimtoenabledesignspaceexplorationthroughause
caseofacceleratingFHEworkloadsonAMDCDNAGPUarchitecture. Leveragingthehigh
flexibilityNaviSimprovides,wedesignanddevelopasetofGPU-basedmicroarchitectural
extensions(GME)totheexistingCDNAarchitecture.
• Wepresentamethodologyforexaminingthekeysimulatordesigndecisionsthatsimulator
developersneedtoconsidertominimizeunnecessarysimulationoverheadandimproveparallel
simulationefficiency. Weintroduceasetofmetricsandmethodstoevaluatethesimulation
speed in terms of various simulated metrics. Based on the profiling results of NaviSim
codebase, we conduct a thorough analysis of the factors that impact the simulation speed
tobetterunderstandthecharacteristicsofdetailed, software-basedarchitecturalsimulation.
Based on the evaluation and analysis, we propose and develop software-based solutions to
improvesimulationspeed.
9

CHAPTER1. INTRODUCTION
1.5 The Organization of this Proposal
This dissertation is organized as follows. Chapter 2 presents background knowledge
associatedwiththisdissertation,including: 1)GPUprogrammingframeworksandarchitectures;2)
thebasicsofarchitecturalsimulators;3)performanceevaluationapproachesandassociatedprofiling
tools; 4) an introduction to fully homomorphic encryption (FHE). Chapter 3 presents prior work
relatedtodevelopingstate-of-the-artGPUsimulatorsandemulators,andextendingGPUarchitecture
toacceleratechallengingworkloads. Chapter4presentourdesignmethodologyandtheNaviSim
GPUsimulator. Chapter5presentsausecasedesignspaceexplorationusingNaviSim. Chapter6
presentsourmethodologyfordesigningandoptimizinghigh-performancesimulators. InChapter7,
weconcludethisdissertationanddiscusssomedirectionsforfuturework.
10

Chapter 2
Background
Inthischapter,weprovideimportantbackgroundthatwillhelpthereaderbetterunderstand
thecontributionsofthisthesis. WebeginbydescribingGPUprogrammingframeworksandGPU
architectures,coveringhowGPUhardwarearchitecturessupportGPUsoftwareprogrammingmodels.
Wethenintroducethebasicsofarchitecturalsimulators. Next,wecoversomeperformanceevaluation
approaches and associated profiling tools. Finally, we introduce fully homomorphic encryption
(FHE), a powerful technique in the privacy-preserving computation domain, which is used as an
exampleapplicationinthiswork.
2.1 GPU Programming Frameworks
A GPU program typically can be separated into the GPU portion and the CPU portion.
TheportionoftheprogramthatexecutesontheGPUiscalledakernel. Theportionoftheprogram
thatrunsontheCPU(i.e.,thehostprogram)launcheskernelsusingvendor-providedruntimeAPIs.
AsetofGPU-controllingAPIsandaGPUprogramminglanguagecomposeaGPUprogramming
framework. CUDA,HIP,andOpenCLarethreeofthemostpopularGPUprogrammingframeworks.
These frameworks share a similar idea of exposing the abstraction of a Single Program Multiple
Data[29](SPMD)modeltoletusersleveragetheunderlyingparallelcomputeenginesinGPUs. We
introducetheprogrammingmodelusingOpenCLterminology,asitcompilesandrunsonNVIDIA,
AMDandIntelGPUs,andotherprogrammingmodels(e.g.,CUDAorHIP)usesimilarsemantics.
Akernelconsistsofanumberofwork-items. Work-itemsaresimilartoCPUthreads;they
executethesameprogramconcurrently,butworkondifferentdata. Work-itemscanbegroupedinto
work-groups. Multiple work-groups compose a kernel in the form of NDRange [10], which can
11

CHAPTER2. BACKGROUND
be arranged in 1D, 2D, or 3D orientation. The work-items in a work-group can be synchronized
usingbarriersandcanshareasmall,butfast,localdatashare(LDS)memory. Allthework-items
of a kernel can access the global memory but there is no way to synchronize all the work-items.
Thisindicatesthatkernel-levelsynchronizationrequiresthehosttostopthekernelandrelaunchthe
kernel.
OnaGPU,asubsetofwork-items(typically32to64)inawork-groupareorganizedin
awavefront. AMDGPUsissueinstructionsatawavefrontgranularity—everytimetheinstruction
scheduler issues one instruction, the ALU repeats the operation for each work-item within the
wavefront. Allthework-itemsinawavefrontruninlock-stepmanner. Withlock-stepexecution,all
thework-itemsinawavefronthavetoexecutethesameinstruction,evenifonlypartofthewavefront
(i.e.,someofthethreads)takesadifferentexecutionpath(i.e.,duetoabranchorthreaddivergence).
ThreaddivergencenegativelyimpactsGPUperformance,asonlypartofthework-itemexecution
willtakeeffect.
GPUframeworksprovideAPIstosupportdatamovementbetweenthehostanddevices.
OpenCLusescommandqueuestomanagetaskssuchasmemorycopiesandcomputekernels. The
commandqueueiscreatedandmanagedbytheGPUdriver,anditservesasacommunicationchannel
betweentheCPUandtheGPU.ItmaintainsasequenceofGPUcommandsthattheGPUwillexecute
inorder. Commandsfromdifferentqueuescanexecuteinparallel.
2.2 Instruction Set Architecture versus Architecture versus Microar-
chitecture
Incomputerarchitecture,InstructionSetArchitecture(ISA),architecture,andmicroarchi-
tecturearecommonlyusedwhendesigninghardwarecomponentsandhardware/softwareinterfaces
to build a computing system. They refer to different aspects of the design and organization of a
computer.
InstructionSetArchitecture(ISA):ISAdefinesthesetofinstructionsthatthecomputing
hardwarecanexecute. Itistheprogrammer’sviewofacomputerandisconcernedwiththeinterface
betweensoftwareandhardware. Itisdefinedbytheinstructionset(language)andoperandlocation
(registerandmemory).
Architecture: Architectureencompassestheentiredesignandorganizationofacomputer
system,includingitscomponents,interconnections,memoryhierarchy,input/outputsystems,and
12

CHAPTER2. BACKGROUND
more. ItisabroaderconceptthatincludesISAasasuperset.
Microarchitecture: Microarchitecturereferstotheimplementationdetailsofthearchitec-
tureatalowerlevel. Itfocusesonthespecificarrangementofregisters,ALUs,memories,andother
building blocks to form a microprocessor. It translates the high-level architecture into a specific,
efficientimplementationbyprovidingorganizationofthedataflowandcontrols,thelogicdesign,
andthephysicalimplementation.
Insummary,ISAdefinestheinstructionsetasaninterfacetotheprogrammer,architecture
encompassestheoverallsystemdesign,andmicroarchitecturedealswiththeinternalimplementation
detailsoftheprocessorforefficientinstructionexecution. Theconceptofarchitectureoverlapswith
ISAandmicroarchitectureanditismostlyused,butISAandmicroarchitecturefocusontheinterface
andphysicalimplementationdetails,respectively.
2.3 GPU Microarchitecture
ComparedwithaCPU,aGPUprovidesmuchhighercomputationalthroughputandmem-
orybandwidth,primarlybasedondifferencesintheorganizationandnumberofdifferentfunctional
units. WeprovideanoverviewofAMDGPUmicroarchitecture,usingCDNAarchitecture[12]asan
example.
TheCDNAarchitecture(seeFigure2.1)adoptsahighlymodulardesignthatincorporatesa
CommandProcessor,ShaderArrays(includingComputeUnitsandL1caches),anon-chipnetwork
connectingthecore-sideL1cachesandthememory-sideL2caches,andDRAM.
2.3.1 KernelDispatching
When a GPU program needs to launch a kernel on a GPU, the program communicates
withthedrivertosendacommandtotheGPU.TheGPUemploysaCommandProcessor(CP)to
processalltheincomingcommandsfromtheCPU,includingmemorycopyingandkernellaunches.
TheCPallocatesanAsynchronousComputeEngine(ACE)toworkasawork-groupdispatcherto
performkerneldispatch. TheACEunitisresponsibleforbreakingdownkernelsintowork-groups
andwavefronts, aswellasassigningthework-groupstotheavailableComputeUnits(CUs). All
work-itemsfromonework-groupmustexecuteonthesameCU,althoughaCUcanconcurrently
processwork-itemsfrommultiplework-groups. Ingeneral,anACEunitcanusedifferentdispatching
algorithms(e.g. roundrobinorgreedy),tofillupeachCUwithanumberofwork-groups. Ingeneral,
13

| CHAPTER2. |     | BACKGROUND |     |     |     |     |     |     |
| --------- | --- | ---------- | --- | --- | --- | --- | --- | --- |
Command Processor
|     |     |     | Shader Engine |     |     | Shader Engine |     |     |
| --- | --- | --- | ------------- | --- | --- | ------------- | --- | --- |
L1I
|              |     | Instruction Cache |              |          |          | Instruction Cache |              | L1I      |
| ------------ | --- | ----------------- | ------------ | -------- | -------- | ----------------- | ------------ | -------- |
|              | L1V |                   | Compute Unit |          |          | Compute Unit      |              | L1V      |
|              | L1V |                   | Compute Unit |          |          | Compute Unit      |              | L1V      |
|              | .   |                   |              | .        |          |                   | .            | .        |
|              | .   |                   |              | .        |          |                   |              |          |
|              |     |                   |              |          |          |                   | .            | .        |
|              | .   |                   |              | .        |          |                   | .            | .        |
|              | L1V |                   | Compute Unit |          |          | Compute Unit      |              | L1V      |
|              | L1S |                   |              |          |          |                   |              | L1S      |
|              |     | Scalar Cache      |              |          |          |                   | Scalar Cache |          |
| Interconnect |     |                   |              |          |          |                   | Interconnect |          |
| L2 Cache     |     | L2 Cache          |              | L2 Cache | L2 Cache |                   | ... ...      | L2 Cache |
... ...
| DRAM |     |     | DRAM       | DRAM                          | DRAM |     |     | DRAM |
| ---- | --- | --- | ---------- | ----------------------------- | ---- | --- | --- | ---- |
|      |     |     | Figure2.1: | TheMicrorchitectureofCDNAGPUs |      |     |     |      |
thenumberofwork-groupsinakernelislargerthanthenumberofCUsonaGPUcanhandleat
the same time. In this case, the ACE fills all the CUs to their maximum capacity first and then
pausesdispatching. OnceanyCUcompletestheexecutionofawork-group,theACEthendispatches
anotherwork-grouptothatCU.Duetothisstaggeredexecutionofwork-groupswithinonekernel,
theGPUlacksthecapabilitytosynchronizeallwork-itemssimultaneously. TherearemultipleACEs
inoneGPU,allowingaGPUtorunmultiplekernelsatthesametime.
2.3.2 ComputeUnit
TheComputeUnit(CU)isresponsibleforinstructionexecutionanddataprocessing. As
showninFigure2.2,aCUisequippedwithascheduler,asetofinstructiondecoders,asetofregister
files,asetofALUs(4SIMDunitsandascalarunit),andmemoryinterfaces. Theschedulerdecides
whichwavefrontcanfetchinstructionsandwhichwavefrontcanissueitsinstructionstobedecoded.
Up to 40 wavefronts can be fetched by the scheduler and stored in wavefront instruction buffers.
There are 6 different execution units, including a branch unit, a scalar unit, a Local Data Share
14

CHAPTER2. BACKGROUND
L1I L1S L1V
Operand Gathering
hcteF
tsnI
10 Wave
Inst Buf
10 Wave
Inst Buf
10 Wave
Inst Buf
10 Wave
Inst Buf
eussI
&
edoceD
To Other CUs
Global
Memory
Register
Access
Write Back
Vec. LDS
VReg SReg SALU Mem Unit
S
-
I
1
M
6
D M
U
a
n
t
i
r
t
ix S
-
I
1
M
6
D M
U
a
n
t
i
r
t
ix S
-
I
1
M
6
D M
U
a
n
t
i
r
t
ix S
-
I
1
M
6
D M
U
a
n
t
i
r
t
ix
Figure2.2: TheMicrorchitectureofCDNAComputeUnit
(LDS) unit, a vector memory unit, four Single-Instruction Multiple-Data (SIMD) units, and four
matrixcoreengines. Thescalarunitalsoservesasascalarmemoryunit,responsibleforexecuting
instructions that manipulate and load data shared by work-items in a wavefront. The LDS unit
executes Local Data Share instructions and interacts with the local memory space. The vector
memoryunitisresponsibleforloadingandstoringmultipledata(64dataelementsinonewavefront
inCDNA)inoneinstruction. EachSIMDunitisresponsibleforexecutingvectorizedfloating-point
instructionsfor10outofthe40wavefrontsmanagedbythescheduler. EachSIMDunitisequipped
with16single-precisionArithmeticLogicUnits(ALUs). Therefore,each64-work-itemwavefront
takes 4 cycles to finish the execution of one instruction. Each matrix engine handles a family of
wavefront-levelinstructions—theMatrixFusedMultiplyAdd(MFMA).MFMAinstructionsperform
mixed-precisionarithmeticandoperateonmatrices.
2.3.3 TheGPUMemoryHierarchy
The CDNA microarchitecture has a two-level cache hierarchy (see Figure2.1). The L1
cachecanbedividedintotheL1scalarcache(mainlyusedforstoringconstantdata,suchaskernel
argumentsandpointers), anL1instructioncache, andanL1vectorcache(awrite-throughcache
thatstoresmostofthedatarequiredbyaCU).TheyareshownasL1S,L1I,andL1Vinthefigure.
Each CU has a dedicated L1 vector cache. CUs in a Shader Array share an L1 scalar cache and
15

CHAPTER2. BACKGROUND
anL1instructioncache. AlltheL1cachesfetchdatafromL2caches(L2sarewrite-backcaches).
Each L2 cache interfaces to a DRAM controller (typically implemented in HBM [30] or GDDR
technology[31]). TheL2cachesandtheDRAMcontrollersarebanked. EachL2cacheservesa
memorybankofacertainaddressrange. Intotal,thereare4to8L2cachesthatcoverthewhole
addressrangeoftheDRAM.TheL1cachesandtheL2cachesareconnectedwithacrossbar. The
crossbardesigncanprovidelow-latencyandhigh-throughputcommunicationchannelsfortheL1
andtheL2caches.
2.4 Architectural Simulators
Architectural simulation plays an important role in the realm of computer architecture
designandresearch,offeringpowerfulwaystomodelandanalyzetheintricatebehaviorofcomputer
systems. Atitscore,architecturalsimulationinvolvesthecreationofsoftware-basedmodelsthat
mimicthebehaviorofhardwarecomponentswithinacomputerarchitecture. Thisenablesresearchers,
engineers,andsystemdesignerstoexperimentwithandevaluatedifferentarchitecturalconfigurations
withouttheneedforphysicalhardware(i.e.,pre-silicon).
2.4.1 MicroarchitectureSimulation
The two main goals of microarchitectural simulation are: i) to model new design ideas
for some parts of a computer system and ii) to enable the evaluation, prediction, and analysis of
thecycle-levelperformance,area,andpower/energy[32]. Toachievethesegoals,software-based
simulationmethodsareused. Eachsimulatorisdevelopedinahigh-levelprogramminglanguage.
Thegoalistomodelthebehaviorofthehardware. Someofthebehaviorofthehardwareismodeled
at a high level of abstraction (i.e., the exact time spent on wires is estimated, though the system
behaves deterministically). The purpose for this abstraction is to simplify the execution of the
simulator(i.e.,noteverycycleismodeledprecisely). Thissimplificationhelpstospeeduptheoverall
runtimeofthesimulation.
Intermsofdesigningasimulator,acriticaldecisiontomakeisselectingtherightlevelof
abstractappropriateforthemodel. Thismayimpactboththespeedofthesimulator(i.e.,themore
detailedasimulator,thegreatertheslowdown)andtheaccuracyofthesimulation(amoredetailed
simulatorshouldproducemoreaccuratesimulationresults).
16

CHAPTER2. BACKGROUND
An architectural simulator is a detailed software modeling tool that fulfills these goals.
Discrete event simulation [33] (DES) is one of the most commonly and widely used methods of
simulationinthecomputerarchitecturefield. Discreteeventsimulationinvolvesmodelingthesystem
asasequenceofeventsthatoccuratdistinctpointsintime,withthestateofthesystemchanging
onlyattheseeventpoints. Keycharacteristicsofdiscreteeventsimulationinclude:
Event basedModeling: Thesimulationmodelisbuiltoneventsthatrepresenttheworkoractions
thatthesimulatorneedstodoinacertainfuturetime. Theseeventsgethandledatspecific
pointsintimeanddrivethesimulationforward.
StateChanges: The state of the system changes only when events occur. Between events, the
systemremainsinastablestate. ThismakesDESparticularlysuitableformodelingsystems
withdynamicandvariablebehaviors.
ConcurrencyandSynchronization: DESallowsforthemodelingofconcurrentprocessesandthe
synchronizationofevents,makingitwell-suitedforsystemswheremultipleentitiesinteractor
competeforresources. TheParallelDiscreteEventSimulation(PDES)technique[34]isused
tosimulatedifferentcomponentsinparallelandconservativelyoroptimisticallysynchronize
theminfixedintervals.
Flexibility: DESinherentlysupportsflexibleandextensibledesigns,giventhatthereareminimal
dependenciesbetweencomponentswithinthesimulator. DESallowscustomizeddefinitions
ofeventsandintegratesseamlesslywithothersoftwaretechniquesthatenhanceadaptability.
Real-TimeSimulation: DESisoftenusedforreal-timesimulation[35],whereeventsaremodeled
to occur at specific, discrete points in time, making it well-suited for modeling computer
systemsoperatingincycles.
2.4.2 ClassificationofSimulators
Simulatorscanbedividedintovarioustypesdependingonthelevelofsimulationdetail,
theinputtothesimulatorandwhethercyclescanbeskipped.
Functional versusTimingSimulators: Functional simulators model thefunctionality
of the target device by recreating the hardware output. Functional simulators can provide basic
programexecutioninformationsuchasinstructioncountsandcachehit/missrates. Incontrast,timing
simulatorsmodelmicroarchitecturaldetails, attemptingtomimictheoperationsastheyoccuron
17

CHAPTER2. BACKGROUND
thehardware. Timingsimulatorskeeptrackofallofthecyclestoprovidedetailedtiming-related
statistics.
Trace-DrivenversusExecution-DrivenSimulators: Trace-drivensimulatorsusetrace
filesastheirinputs. Thesefilescontainrecordedstreamsofinstructionsand/oraddressescaptured
duringapreviousexecutiononthetargethardware. Thesimulatorreadsatraceofinstructionsor
addressesthatwerepreviouslyexecuted,sothereisnoneedtodecodetheISAofthetarget,which
makesthemrelativelysimple. Trace-drivensimulatorsareconstrainedtostudyingspecificscenarios
and conditions, and are unable to support new designs that rely on dynamic information. In an
execution-drivensimulator,instructionsarefirstdecodedandexecutedonthesimulatedmachine
directly. Execution-drivensimulatorsprovidegreaterflexibility,asthereisnoneedtoexecutethe
applicationonrealhardwarefirst. Execution-drivensimulationsupportsmoredetailedmodelingof
complexinteractions,insteadofreplayingpre-recordedsystembehavior. However,execution-driven
simulatorsaremorecomplexcomparedtotrace-drivensimulators(e.g.,thesesimulatorsmayrequire
tomodelaninstructiondecoderandinstructioncaches)andmaytakelongertosimulate.
Cycle-basedversusEvent-basedSimulators: Inabroadersense,acycle-basedsimulator
is considered a DES. It is a simplified DES having events occur on every clock cycle. The main
differencebetweencycle-basedandevent-basedsimulatorsliesinthefactthatcycle-basedsimulators
cannotskipanycycle,evenifthereisnoworktodointhatcycle,whichincreasessimulationtime
unnecessarily. Manysimulatorscombinethesetwobymodelingsomepartsofthesimulatedtarget
onacycle-by-cyclebasisandothersonanevent-drivenbasis.
2.5 Performance Evaluation Tools
Measuringandanalyzingtheperformance, specificallytiminginformation, ofanappli-
cation involves collecting performance counter information from various hardware and software
components. Various tools and techniques are available to support performance tuning, includ-
ingprofilingtoolsprovidedbytheoperatingsystem(e.g., perf[36]andgprof[37]forLinuxand
WindowsPerformanceToolkit(WPT)[38]forWindows),specificvendors(e.g.,Intel,AMD),and
cross-platformlow-levellibraries. Sometoolsfocusonuser-levelcodeinspecificlanguages,suchas
gprof[37]forCandFortran,pprof[39]forGolang.
18

CHAPTER2. BACKGROUND
2.5.1 perf
perf[36]isaLinuxcommand-linetoolthatoperatesatthesystemlevel. Itprovidesawide
rangeofperformancemonitoringfeatures,includingCPUutilization,memoryaccesspatterns,cache
misses, and other hardware-level events. Perf is not limited to profiling only user-level code, but
extendstosystem-wideanalysis,allowingyoutoinvestigatetheperformanceoftheentiresystem.
It is often used for analyzing low-level system performance, kernel-level events, and hardware
performancecounters. perfcomespre-installedonmanyLinuxdistributions.
2.5.2 PAPI
PerformanceApplicationProgrammingInterface[40](PAPI),isalibraryandsetoftools
designed to provide a consistent interface and methodology for collecting performance counter
informationfromtheCPUandotherhardwarecomponents. Itallowsdeveloperstoaccesslow-level
hardwareperformancecountersinaplatform-independentway,makingiteasiertogatherdetailed
performancemetricsacrossdifferentarchitectures. PAPIincludesasetofAPIsthatdeveloperscan
usetointegrateperformancemonitoringcapabilitiesintotheirapplications. TheseAPIssupportC
andFortranprogramminglanguages.
2.5.3 pprof
pprof is a profiling tool associated with the Go programming language. It is used for
collecting and analyzing profiling data from Go applications. pprof is a part of the Go standard
library and provides built-in support for profiling Go applications. It allows developers to gather
insights into the runtime behavior of their Go code. Pprof supports CPU timing and memory
profiling. IthelpsidentifythefunctionsandcodepathswheremostoftheCPUtimeisspentduring
the execution of a Go program. It also helps developers understand memory allocation patterns,
identifymemoryleaks,andoptimizememoryusageinGoprograms.
2.6 Fully Homomorphic Encryption
Fullyhomomorphicencryption(FHE)isagroundbreakingcryptographictechniquethat
emergedasasolutiontoalongstandingchallengeinthefieldofsecuredataprocessingandprivacy.
Thisformofencryptionallowscomputationstobeperformedonencrypteddatawithouteverrevealing
theunderlyinginformation(SeeFigure2.3). Thisfurtherpavesthewayforawiderangeofpractical
19

CHAPTER2. BACKGROUND
applications,fromsecureoutsourcedcomputationandprivatedataanalysistoprivacy-preserving
machinelearning.
Internet
Safe Hostile Untrusted
Environment Environment Environment
Trust
Barrier
Message Encrypt Cipher-text Cipher-text
Message Decrypt Cipher-text Cipher-text
srevreS
GME
GPUs
Figure2.3: FHEoffersasafeguardagainstonlineeavesdroppersaswellasuntrustedcloudservices
byallowingdirectcomputationonencrypteddata.
Theideaofhomomorphicencryptionwasproposedin1978byWhitfieldDiffieandMartin
Hellman. A few earlier homomorphic encryption schemes are partial homomorphic encryption
schemes, because they can support only limited operations and have limited usability. Paillier’s
cryptosystem[41]couldnotsupportmultiplicationoperations,onlyadditionoperationsaresupported,
whiletheRSAcryptosystem[42]onlysupportsmultiplication,butnotaddition. Thenin2009,Gentry
presented,inhisseminalpaper[43],anewapproachtoFHEthatovercamesomeofthekeychallenges
associatedwithearlierproposals. Gentry’sapproachuseslattice-basedcryptographyandenablesboth
additionandmultiplicationoperations. Theinitialversionofthehomomorphicencryptionscheme
was incomplete because, while it could support any number of additions, it could only support a
limitednumberofmultiplications. Thenlateron,Gentrymadesomeimprovementsbyintroducing
anoperationknownasbootstrapping. Thisimprovementelevatedhomomorphicencryptiontoanew
level. ThisnewapproachtoFHEwascalledfullyhomomorphicencryptionbecauseitcouldsupport
anynumberofadditionsandanynumberofmultiplications,turningthetheoreticalconceptoffully
homomorphicencryptionintoapracticalreality.
2.6.1 TheCKKSFHEScheme
Inthisthesis,wewillutilizedtheCheon-Kim-Kim-Song(CKKS)FHEschemetodemon-
stratetheflexibilityofoursimulationinfrastructuretobeadaptedtoacceleratethiscomplexscheme.
CKKS can support a wide range of privacy-preserving applications by allowing operations on
20

CHAPTER2. BACKGROUND
Figure2.4: TheoperationflowoftheCKKSscheme.
floating-point data. Figure 2.4 provides a high-level view of CKKS showing the flow of the op-
erationsperformedinthescheme. Webeginwithamessagem. Themessagespaceisdefinedas
a space of complex numbers. Basically, this is a vector of complex values on which we want to
perform certain computation. The first step is to encode this message to generate what we call
plaintext. Theplaintextspaceisconstrainedbyanirreduciblepolynomialandthevaluesarethe
coefficientsofthisplaintextpolynomial
Z[x].
Afterencoding,weencrypttheplaintext,whichnow
is represented as a single polynomial, into a pair of polynomials to generate the ciphertext using
apublickey. Thecoefficientsoftheencryptedpolynomialaredefinedas Z [x]. Weintroducethe
q
polynomial modulus Q to constrain the size of each coefficient. We can examine the ciphertext
further to understand its content. Basically, the ciphertext consists of a pair of elements in the
polynomial ring R Q =Z Q [x]/(xN+1). Each element of this ring is a polynomial ∑ i N = − 0 1a i xi with
“degree-bound”N−1andcoefficientsa in Z . Foramessagem∈Cn,wedenoteitsencryptionas
i Q
m =(A ,B )whereA andB arethetwopolynomialsthatcomprisetheciphertext. Afterthe
m m m m
(cid:74) (cid:75)
messageisencodedandencrypted,wecanperformhomomorphicencryptionoperations(denotedas
f). Thenwegetaresultinciphertext,representedasc′. Weneedtodecryptittoplaintextandthen
decodeitbacktoaresultmessagem′. ThebasicdatatypeinCKKSisavectorwithlengthN inboth
plaintextandciphertext,andeachelementinthevectorisofsizelogQbits. N andQarethemain
parameters,whichdefinethesizeoftheciphertext,andalsogovernthesizeoftheworkingdataset
thatisrequiredtobepresentintheon-chipmemory.
21

CHAPTER2. BACKGROUND
2.6.2 BuildingBlocksoftheCKKSFHEScheme
SinceCKKSschemedealswiththeLearningwithErrors(LWE)[44]problem,oneofthe
challengesishowtohandletheproblemthattheerrorintroducedbyhomomorphicencryption(HE)
operationsgrowsasweperformasequenceoftheseoperations. Forexample,theadditionoperation
doubles the noise, and multiplication squares the noise. The noise keeps growing if we perform
multiplelevelsofadditionormultiplicationoperationsthenatonepointtheerrorissolargethatwe
arenolongerabletodecryptit. Therearetwosolutionstothisproblem. Oneiswecanchooselarge
enoughparameterssothatwehavealargespaceforthenoisetogrow. Anothersolutionisthatwe
canperformabootstrappingoperation.
Choosing a large enough value for logQ will help us preserve the underlying message
because it allows more space for the error to grow, without losing precision. Also, there is a
correlationbetweenthevalueforlogQandthevalueofN. AlargeenoughlogQvaluerequiresa
largeenoughvalueofN sothatwehaveenoughsecuritybitsfortheunderlyingencryptionscheme.
For128-bitsecurity, typicalvaluesofN rangefrom216 to217, andlogQvaluesrangefrom1700
to 2200 bits for practical purposes. These large sizes of N and logQ are required to maintain the
securityoftheunderlyingRing-LearningwithErrorsassumption[44].
However,therearenocommerciallyavailablecomputesystemsthathaveahundred-bit
wideorthousand-bitwideALU,whicharenecessarytoprocesstheselargecoefficients. Acommon
approachforimplementingtheCKKSschemeonhardwarewithamuchsmallerwordlengthisto
chooseQtobeaproductofdistinctword-sizedprimesq ,...,q . Then Z canbeidentifiedwith
1 ℓ Q
the“productring”∏ l
i=1
Z
qi
viatheChineseRemainderTheorem[45]. Inpractice,thismeansthat
the elements of Z can be represented as an ℓ-tuple (x ,...,x ), where x ∈Z for each i. This
Q 1 ℓ i qi
representationofelementsin Z isreferredtoastheResidueNumberSystem(RNS)andiscommonly
Q
referredtoasthelimbsoftheciphertext.
TheCKKSFHEschemesupportsfloating-pointmessages,soallencodedmessagesmust
includeascalingfactor∆withthepurposeofreducingtheprecisionloss. CKKSdealswithrealand
complexnumbers,whicharecontinuousandhaveinfiniteprecision. Incontrast,digitalcomputers
workwithfiniteprecision. Thescalefactorisusedtomaprealorcomplexnumberstoafiniteand
discretesetofvalues. Itiscrucialtoconstrainerrorgrowthbecauseithelpsensurethatciphertexts
remain proportional to the scale factor, while the error remains the same. This scaling factor is
typicallythesizeofoneofthelimbsoftheciphertext,soitsvalue∆istypicallyequivalenttologq.
Aswecansee,sinceciphertextsarerepresentedbypolynomials,polynomialmultiplication
22

CHAPTER2. BACKGROUND
phi W
Polynomial1 q W-1 1/N phi-1
NWC NTT
Pointwise Scalar Polynomial3
iNTT iNWC
Multiplication Mult.
Polynomial2
NWC NTT
Figure2.5: PolynomialMultiplication.
becomes one of the main operations in the CKKS scheme, and more generally, in lattice-based
cryptography. Polynomialmultiplicationisaveryexpensiveoperationandinordertoenablefast
polynomial multiplication, by default, we represent polynomials as a series of N evaluations at
fixedrootsofunity[46]. Werefertothispolynomialrepresentationastheevaluationrepresentation.
Movingbetweenthetwopolynomialrepresentationsrequiresanumber-theoretictransform(NTT)
or inverse NTT, which is the finite field version of the fast Fourier transform (FFT)[43]. After
performing NTT on polynomials we can perform point-wise multiplication of the corresponding
coefficientsofthetwopolynomials. ThisallowspolynomialmultiplicationtooccurinO(N)time
insteadofO(N2)time. Figure2.5showsthestepsinvolvedinpolynomialmultiplication. Itconsists
ofperforminganegativewrappedconvolutionoperation(NWC),followedbyNTToperations,which
isthenfollowedbyapoint-wisemultiplication. NWCisamultiplicationwith phiwhichisgenerated
from the root of unity. ω and phi are N and 2N root of unity respectively. After the point-wise
multiplication,wecanperformaninverseNTT,thendoascalarmultiplicationwith1/N,andthen
performaninversenegativewrappedconvolution(iNWC).
We list all the building blocks in the CKKS scheme in Table 2.1. All of the operations
thatformthebuildingblocksoftheCKKSschemeinvolvepolynomialadditionandmultiplication,
whichfurtherreducetoscalarmodularadditionsandscalarmodularmultiplications. TheScalarAdd
andScalarMultarethetwomostbasicbuildingblocksthataddandmultiplyascalarconstanttoa
ciphertext. PolyAddandPolyMultaddandmultiplyaplaintextpolynomialtoaciphertext. Wedefine
separateScalarAddandScalarMultoperations(inadditiontoPolyAddandPolyMult)becausethe
scalarconstantvaluescanbefetcheddirectlyfromtheregisterfile,whichcanhelpsaveexpensive
mainmemoryaccesses. NotethatthePolyMultisfollowedbyanHERescaleoperationtorestore
thescaleofaciphertextto∆fromscale∆2. Aswementionedbefore,whenmultiplyingmessages
together,thisscalingfactorgrowsaswell. Thescalingfactormustbeshrunkdowninordertoavoid
overflowingtheciphertextcoefficientmodulus.
23

| CHAPTER2.  | BACKGROUND        |           |                               |          |     |     |     |                          |
| ---------- | ----------------- | --------- | ----------------------------- | -------- | --- | --- | --- | ------------------------ |
|            |                   | Table2.1: | TheHEbuildingblocksusingCKKS. |          |     |     |     |                          |
| Block      |                   |           | Computation                   |          |     |     |     | Description              |
| ScalarAdd( | m ,c)             |           | m+c                           | =(B +c,A | )   |     |     | Addascalarctoaciphertext |
|            |                   |           |                               | m        | m   |     |     |                          |
|            | (cid:74) (cid:75) |           | (cid:74)                      | (cid:75) |     |     |     |                          |
where,cisalength-N vectorwithevery
elementc
ScalarMult( m ,c) m·c =(B ·c,A ·c) Multiplyascalarbyaciphertext
m m
|          | (cid:74) (cid:75) |     | (cid:74) | (cid:75) |       |     |     |                            |
| -------- | ----------------- | --- | -------- | -------- | ----- | --- | --- | -------------------------- |
| PolyAdd( | m ,P )            |     | m+m′     | =(B      | +P ,A | )   |     |                            |
|          | m′                |     |          | m        | m′    | m   |     | Addanunencryptedpolynomial |
|          | (cid:74) (cid:75) |     | (cid:74) | (cid:75) |       |     |     |                            |
toaciphertext
| PolyMult( | m ,P )            |     | m·m′     | =(B ∗P   | ,A   | ∗P ) |     |                          |
| --------- | ----------------- | --- | -------- | -------- | ---- | ---- | --- | ------------------------ |
|           | m′                |     |          | m        | m′ m | m′   |     | Multiplyinganunencrypted |
|           | (cid:74) (cid:75) |     | (cid:74) | (cid:75) |      |      |     |                          |
polynomialwithaciphertext
|        | m′                                  |     | m+m′     |          |          |         |     |                   |
| ------ | ----------------------------------- | --- | -------- | -------- | -------- | ------- | --- | ----------------- |
| HEAdd( | m , )                               |     |          | =(B m    | +B m′ ,A | m +A m′ | )   | Addtwociphertexts |
|        | (cid:74) (cid:75) (cid:74) (cid:75) |     | (cid:74) | (cid:75) |          |         |     |                   |
|        | m′                                  |     | m·m′     |          |          |         |     |                   |
HEMult( m , ,evk ) =KeySwitch(A ∗A m′ ,evk )+ Multiplytwociphertexts
|     |                                     | mult |          |          |       | m    | mult |     |
| --- | ----------------------------------- | ---- | -------- | -------- | ----- | ---- | ---- | --- |
|     | (cid:74) (cid:75) (cid:74) (cid:75) |      | (cid:74) | (cid:75) |       |      |      |     |
|     |                                     |      | (B ∗B    | ,A ∗B    | +A    | ∗B ) |      |     |
|     |                                     |      | m        | m′ m     | m′ m′ | m    |      |     |
|     |                                     | (r)  |          |          |       |      | (r)  |     |
HERotate( m ,r,evk ) m≪r =KeySwitch(ψ (A ),evk )+ Circularrotateelementsleftbyr
|     |                   | rot |          |          |     | r m | rot |                            |
| --- | ----------------- | --- | -------- | -------- | --- | --- | --- | -------------------------- |
|     | (cid:74) (cid:75) |     | (cid:74) | (cid:75) |     |     |     |                            |
|     |                   |     | (ψ (B    | ),0)     |     |     |     |                            |
|     |                   |     | r        | m        |     |     |     | slots,ψ r isanautomorphism |
performed
| HERescale( | m )               |     | ∆−1·m    | =(∆−1B   | ,∆−1A | )   |     |                              |
| ---------- | ----------------- | --- | -------- | -------- | ----- | --- | --- | ---------------------------- |
|            |                   |     |          |          | m     | m   |     | Restorethescaleofaciphertext |
|            | (cid:74) (cid:75) |     | (cid:74) | (cid:75) |       |     |     |                              |
fromscale∆2 backto∆
24

CHAPTER2. BACKGROUND
TheHEAddoperationisstraightforwardandaddsthecorrespondingpolynomialswithin
thetwociphertexts. However,theHEMultandHERotateoperationsarecomputationallyexpensive,
astheyperformakeyswitchoperationafterthemultiplicationandautomorphoperations,respectively.
In both the HEMult and HERotate implementations, there is an intermediate ciphertext with a
decryptionkeythatdiffersfromthedecryptionkeyoftheinputciphertexts. Inordertotransform
thisnewdecryptionkeybacktotheoriginaldecryptionkey,weperformakeyswitchoperation. This
(r)
operationtakesinaswitchingkey(eitherevk orevk ),aciphertextandaciphertext m thatis
mult rot s
(cid:74) (cid:75)
decryptableunderasecretkeys. Theoutputofthekeyswitchoperationisaciphertext m s′ that
(cid:74) (cid:75)
encryptsthesamemessage,butisdecryptableunderadifferentkeys′.
Keyswitchoperationsitselfarehighlycomputationallyexpensiveanditinvolvesasequence
ofsub-operations. Toincurminimalnoisegrowthduringthekeyswitchoperation,thekeyswitch
operation requires that a polynomial be split into dnum digits, then raise the modulus before
multiplying with the switching key, followed by a modulus down operation. The modulus raise
and down operations operate on the coefficient representation of the polynomial, requiring us to
performexpensiveNTTandiNTTconversions. Moreover,theswitchingkeysarethesamesizeas
theciphertextitself,requiringustofetchalargeamountofdataatoncetomultiplytheswitching
keyswiththeciphertext. Finally,thereexistsanoperationknownasbootstrapping[43]thatneeds
tobeperformedfrequentlytode-noisetheciphertext. Thisbootstrappingoperationisasequence
ofthebasicbuildingblocksintheCKKSscheme,meaningthatitsuffersfromthesamecompute
and memory bottlenecks that exist in these building blocks, making it one of the most expensive
operations. WelisttheparametersthatdefinetheCKKSFHEschemeinTable2.2.
25

CHAPTER2. BACKGROUND
| Table2.2: | TheCKKSparametersanddescriptions. |      |     |     |
| --------- | --------------------------------- | ---- | --- | --- |
| Param     | Description                       |      |     |     |
| N         | Polynomialdegree-bound            |      |     |     |
| n         |                                   | n≤ N |     |     |
Lengthofthemessage.
2
| Q   | Polynomialmodulus                 |     |     |     |
| --- | --------------------------------- | --- | --- | --- |
| L   | Maximumnumberoflimbsinaciphertext |     |     |     |
C
|      | Theset{q ,q                           | ,...,q L }ofprimefactorsofQ |     |     |
| ---- | ------------------------------------- | --------------------------- | --- | --- |
|      | 0                                     | 1                           |     |     |
| ℓ    | Numberoflimbs,numberoffactorsinQ;     |                             |     |     |
| dnum | Numberofdigitsintheswitchingkey       |                             |     |     |
| α    | Numberoflimbsthatcompriseasingledigit |                             |     |     |
=⌈L+1⌉
inthekey-switchingdecompositionα
dnum
| P       | Productofextensionlimbsaddedfor    |                         |     |     |
| ------- | ---------------------------------- | ----------------------- | --- | --- |
|         | raisedmodulus.                     | Totalextensionlimbs=α+1 |     |     |
| fftIter | Multiplicativedepthofbootstrapping |                         |     |     |
lineartransform
| ∆   | Scalemultipliedduringencryption |     |     |     |
| --- | ------------------------------- | --- | --- | --- |
| m   | Amessagevectorofnslots          |     |     |     |
m
Ciphertextencryptingamessage
(cid:74) (cid:75)
| A   | Arandomlysampledpolynomialfrommessagem |     |     |     |
| --- | -------------------------------------- | --- | --- | --- |
m
| P   | Encryptedmessageasapolynomial |     |     |     |
| --- | ----------------------------- | --- | --- | --- |
| P m | Polynomialencryptingmessagem  |     |     |     |
| [P] | q-limbofP                     |     |     |     |
| qi  | i                             |     |     |     |
| evk | Evaluationkey                 |     |     |     |
(r)
| evk | EvaluationkeyforHE-Rotateblockwith |     |     |     |
| --- | ---------------------------------- | --- | --- | --- |
rot
(r)rotations
| Table2.3: | PracticalparametersforourFHEoperations. |        |         |     |
| --------- | --------------------------------------- | ------ | ------- | --- |
| log(q) N  | logQ L                                  | L dnum | fftIter | λ   |
boot
| 54 216 | 1728 23 | 17 3 | 4   | 128 |
| ------ | ------- | ---- | --- | --- |
26

Chapter 3
Related Work
Inthischapter,wereviewpriorworkrelatedtotheresearchareascoveredinthisdissertation.
Wefirstbeginbydiscussingthestate-of-the-artinGPUsimulatorsandemulators. Next,wediscuss
related work on extending GPU architecture to accelerate FHE workloads. There is no directly
relatedworkofextendingGPUarchitectureforFHEbutthemostcloselyrelatedworkfocuseson
DeepLearning(DL)workloads[47].
3.1 GPU Simulators
GPU simulators have been critical infrastructures that enable GPU architecture design
and quantitative evaluation. To date, GPU architecture researchers have devoted major efforts to
developingGPUsimulatorsandemulatorstosupportarchitecturalexplorationinbothindustryand
academia.
3.1.1 Commercially-developSimulators
GPUmanufacturingcompanies,suchasNVIDIAandAMD,havedevelopedtheirown
internalsimulationtools. NVArchSim(NVAS)[28],developedbyNVIDIA,isanadvancedsimulation
frameworkfortestingandrefiningtheirGPUdesigns. NVIDIA’sArchitecturalSimulator(NVArch-
SimorNVAS)focusesonbuildingtrustworthysimulationplatformswithdifferentlevelsoffidelity
toachievehighersimulationspeed. NVASisahybridtrace-drivenandexecution-drivenmulti-GPU
simulator. NVAS is used to evaluate HPC and ML workloads with much higher performance as
comparedtomostpubliclyavailableGPUsimulators.
27

CHAPTER3. RELATEDWORK
AMD has developed their own GPU simulator based on the gem5 simulator [48]. The
AMDgem5GPUmodel[49],whosemaintargetistheAMDAPUdevices,isacomponentaddedto
thegem5simulator. ThenewcomponentisdedicatedtotheAMDGCN3architecture. TheAMD
gem5GPUsimulatorisanexecution-drivensingle-GPUsimulator.
3.1.2 AcademicResearchSimulators
DevelopingGPUresearch-relatedtoolshasbeenafocalpointsincethetremendousgrowth
inusinggraphicsprocessorsforcompute. EarliertoolsincludeBarra[50]andGPUOcelot[51,52],
whichprovidefunctionalGPUemulationsupport,thoughdonotsupporttimingmodels. Overthepast
decade,manypubliclyavailableGPUsimulatorshavebeendevelopedtocapturenotonlyfunctional,
butalsotimingbehavior,ofvariousgenerationsofGPUs. WepresenttheseGPUsimulatorsbased
onthetimingoftheirGPUtarget. ThenweintroduceseveralparallelGPUsimulators.
ForNVIDIAGPUs,GPGPUSim[53]isaGPUsimulatorthatmodelstheNVIDIAFermi
architectureandsupportstheCUDA4.0runtimeAPI.GPGPUSimsupportsvirtualISAPTX,an
intermediatelanguagebetweentheCUDAprogramminglanguageandasubsetofanoldermachine
ISASource and ASSembly(SASS).GPGPUSimhasbeenextendedtosupportadditionalGPU
features. MAFIA[54]extendsGPGPUSimtosupportconcurrentkernelexecutionandMOSAIC[55]
extendsMAFIAsimulatortosupportvirtual-to-physicaladdresstranslationandaunifiedCPU-GPU
addressspace. Accel-Sim[56]isanotherworkextendedfromGPGPUSim. Itisatrace-basedGPU
simulatorwithaflexiblefrontendenablingthesupportofmodernSASSmachineISAs. Accel-Sim
includesatrace-generationtool(builtusingtheNVBit[57]binaryinstrumentationtool)toproduce
machineISAinstructiontracesfromanyCUDAbinary. Theyalsomakeextensivemodificationsto
GPGPUSim’sperformancemodeltomodelfeaturesintroducedinmoremodernGPUs. Accel-Sim
alsoaddssupportfortheNVIDIAKepler,Pascal,Volta,andTuringarchitectures.
ForAMDGPUs,Multi2Sim[58]isoneoftheearlysimulationtoolsthatsupportedAMD
GPU simulation. Multi2Sim is a simulator infrastructure capable of emulating various CPU and
GPUarchitectures. Multi2SimsupportsdetailedtimingsimulationforCPUsoperatingontheX86
instructionset. IntermsofGPUs,Multi2SiminitiallymodeledtheAMDEvergreenandSouthern
IslandsseriesGPUs. Notably,ithasexpandeditssupporttoincludethemodelingofNVIDIAKepler
GPUs[59]. MGPUSim[20]isahigh-performanceparallelGPUsimulatortargetingtheAMDGCN3
architecture. MGPUSim can run either emulation mode or timing simulation mode and provides
nativelanguage-levelsupportformulti-threadedprogramming. Moreover,MGPUSimprovidesa
28

CHAPTER3. RELATEDWORK
flexiblemulti-GPUsimulationframeworksupportingbothunifiedanddiscretemulti-GPUmodels.
ForIntelGPUs,theMacsim[60]isasimulatorframeworkcapableofsimulatingavariety
ofarchitecturesincludingx86,ARM64,NVIDIAPTXandIntelGENGPUarchitecture. Itmodels
thedetailedIntel’sintegratedGPUmicroarchitecture.
The parallel execution model of the GPU leads to a large number of units concurrently
executing independently, which produces opportunities to use multiple CPU threads to simulate
GPU execution. Barra [50] mainly focuses on parallel functional emulation. GPUTejas [61] is a
Java-based, trace-driven, parallel architectural simulator that can achieve high performance and
scalability. Leeetal. proposedaparallelsimulatorframework[62,63],whichmodifiesGPGPUSim
andonlysynchronizeswhentheprocessoraccessesthememorysystem. MGPUSimalsosupports
parallel simulation adopting a conservative parallel simulation approach which requires global
synchronization after each cycle. Hughes et al. integrated GPGPUSim into a node mode in the
SST[64]frameworktopresentaparallelizable,multi-node,GPUsimulator[65].
3.2 Extending GPU Architectures to Accelerate FHE
3.2.1 Extending/ModifyingCurrentGPUArchitectures
Researchers have been trying to address challenges in workload-related problems (i.e.,
accelerationofspecificapplicationsandconcurrentexecutionsofmultipleapplications)bydesigning
GPUarchitectureextensions.
Marangozetal.[66]proposeanextensionoftheGPUmemoryarchitecturetoeffectively
supportbandwidthreservationandsharingwithlittleoverhead. Theyaddresstheinterferenceproblem
betweentheco-runningapplicationswhichincreasesnon-deterministictimingbehaviorsbydesigning
aGPUmemorybandwidthreservationschemeandimplementingnewhardwarecomponentsina
GPUsimulatorbuiltuponGPGPUSim.
Fu et al. [47] propose a Composable On-PAckage GPU (COPA-GPU) architecture to
accommodate architectural requirements between the FP32 (or larger)-based HPC and FP16 (or
smaller)-based DL workloads. By leveraging multi-chip module (MCM) [67] integration, they
proposecomposable,semi-specializeddesignsthatprovidehighlevelsofGPUdesignreuseacross
theHPCandDLdomains. Thepaperpresentsthemicroarchitectureoftheproposedmemorysystems,
costandbenefitsoftheCOPA-GPU.TheyevaluatetheCOPA-GPUdesignusingamodifiedversion
ofNVIDIA’sNVArchSim.
29

CHAPTER3. RELATEDWORK
Braaketal.[68]proposeR-GPU,areconfigurableGPUarchitecturewithcommunicating
cores. R-GPU,asanextensiontoaGPU,enablescommunicationbetweenthecoresinaStreaming
Multiprocessor(SM,similartoCU)byaddingconfigurationregisters(CR)andacommunication
networkwithFIFObuffers. Thepaperintroducesaprogrammingmodeltoassisttheprogrammer
inprogrammingtheR-GPUarchitectureandprovidesanestimationofthehardwarecostsinterms
ofareaandpowerconsumption. TheideaofconfiguringcoresinanSMinanetworkissimilarto
thefirstmicroarchitecturalextensionpresentedinGME[69],whichbuildsaCU/SM-sidenetwork-
on-chipthatinterconnectsCU/SMstosupportdatamovementbetweenCU/SMs. Theyvalidatethe
performanceimprovementsofR-GPUinGPGPUSim.
3.2.2 AcceleratingFHE
SincethemajorchallengeindeployingFHEinreal-worldsystemsisovercomingthehigh
computationalcosts,manyprioreffortshavebeeninvestedinbuildingefficientFHEimplementations
includinglibrariesandframeworks.
For GPU implementation acceleration, Privft [70] provides a GPU implementation of
theCheon-Kim-Kim-Song(CKKS)FHEschemethatshows1to2ordersofmagnitudespeedup.
TheworkbyBadawietal.[71]aimstoaccelerateFHEusingNVIDIAGPUs. 100x[72]aimsat
speedinguptheend-to-endCKKSschemeincludingbootstrapping. Theyutilizememory-centric
optimizationssuchaskernelfusionandreorderingprimaryfunctionstooptimizeoff-chipmemory
transactions.
SeveralpriorworksaccelerateFHEbydevelopingdomain-specificcustomaccelerators.
HEAX[73]isanFPGA-basedacceleratorthatonlyspeedsupCKKSencryptedmultiplication,with
theremainderoffloadedtothehostprocessor. Kimetal.[74]proposeanFPGA-basedimplementation
for HE operations with the bootstrappable parameters. It only focuses on optimizing Number
Theoretic Transform (NTT), a fundamental operation in lattice-based cryptography, to support
the bootstrapping process in FHE. FAB [75] is the most recent work and it proposes an FPGA-
basedacceleratorforbootstrappableFHEforasecureandpracticalparameterset. FABachieves
performancecomparabletothepreviousGPUimplementation,100x[72]. However,mostofthese
priorproposalseitherdonotcoverallHEprimitives,oronlysupportasmallerparametersetthat
allow computation up to a multiplicative depth of 10, which is insufficient for bootstrapping or
applicationssuchaslogisticregressiontraining.
AlthoughFPGAsshowgreatpotentialforacceleratingFHEworkloads,theyarelimited
30

CHAPTER3. RELATEDWORK
to run at low operating frequencies and provide limited compute resources. Researchers have
developedseveralASIC-acceleratorstoaccelerateFHE.F1[76]proposesanASICimplementation
thatmakesuseofsmallN andQvalues,implementingonlyasingle-slotbootstrapping. BTS[77]
isthefirstASICproposaldemonstratingtheperformanceofafully-packedCKKSbootstrapping.
CraterLake[78]andARK[79]furtherenhancethepackedCKKSbootstrappingperformanceand
demonstrateseveralordersofperformanceimprovementacrossvariousworkloads.
AlthoughASICimplementationsdeliverseveralordersofperformanceimprovement,they
are tailored for a single workload. With customized designs, they can adopt the latest memory
technology. For example, ARK [79] integrates the latest HBM3 technology and CraterLake [78]
incorporatesextraphysicallayerstofacilitatecommunicationbetweenDRAMandon-chipmemory.
Hence,theyareabletoenhancetheavailablebandwidthforFHEworkloads.
31

Chapter 4
NaviSim
The development and availability of architectural simulators has enable industry and
researchtoaggressivelyperformdesignspaceexploration. Architecturalsimulationtoolsarevital
astheyprovidequantitativefeedbackondesignoptimizationandcanaidinpre-siliconverification.
AwiderangeofsimulatorshavebeendevelopedtoexplorethedesignspaceofCPUs,GPUs,and
customeraccelerators. Inthisthesis,weattempttocharacterizetheprocessofdevelopingasimulator,
utilizinganumberofdesigngoals. Togroundourdiscussioninreality,wewillcoverNaviSim,a
novelGPUsimulatorthatmodelstheAMDRDNAandCDNAarchitectures. WhileNaviSimhas
beenusedtoevaluatetheperformanceofAMDCDNAGPUs,wewillleverageheretoexplorethe
differentsimulatormetricsdevelopedinthisthesis.
One of the key motivations for carefully considering how best to architect a new GPU
simulatoristherapidrateofinnovationoccurringinGPUarchitectures. GPUvendorshaveintroduced
next-generationarchitecturesandproductsatanunprecedentedpace. ExistingGPUsimulatorsare
typicallydesignedtosimulateaspecificinstructionsetarchitecture,andarecustomizedtocapture
a very specific GPU microarchitecture. This rigid design makes it difficult to adapt an existing
simulatortosupportfuturearchitectures. Thiswillimpactthearchitecturecommunity’sabilityto
properlyassessfuturearchitectures.
Over the past five years, AMD announced two brand-new architectures — RDNA and
CDNA,movingawayfromtheGCNarchitecturewhichhasbeenthedominantGPUarchitecturefor
AMDforoveradecade. OneuniquegoalofNaviSimistomodelthetwonewarchitecturesalong
withGCN3andcomparethedifferencesbetweenISAsandmicroarchitectures.
32

CHAPTER4. NAVISIM
4.1 Design Goals
An overall design goal of NaviSim is to achieve adaptability. We design NaviSim to
anticipatemultipleISAsandmicroarchitectures,attemptingtoguardagainstthesimulatorbecoming
quicklyoutdatedduetothepaceofdevelopmentinGPUtechnology. Webreakdownadaptability
furtherintotwoqualities: i)modularityandii)extensibility.
Highlymodulardesign. WeseparatetheISAemulationfromthedetailedtimingsimu-
lationbyadoptingastand-aloneinstructionemulator. Asanexecution-drivensimulator,NaviSim
recreatestheexecutionresultsofGPUinstructionsduringsimulationwiththehelpofaninstruction
emulatorforindividualISAs. Ourdesignonlydefinesahigh-levelinterfaceoftheemulatorandthis
allowsuserstoeasilyswaptheinstructionemulatorbeingused.
Inaddition,webreakdownouroveralldesignintoasetofisolatedcomponents,which
serve as the basic simulation elements of the framework. Isolation here means each component
should be a self-sufficient entity that maintains all the state required to perform tasks. More
specifically,acomponentshouldnotdirectlyaccessfieldsofothercomponentsorcallfunctionsof
othercomponents. Thisisaformofdesignencapsulation,akintohowobject-orientedprogramming
languageencapsulatedataandoptionsinaclass. Moreoverinourdesign,asimulatorcomponent
cannotscheduleaneventforanothercomponent.
EaseofDesignExtensions. Whentheresearchcommunityusesanarchitecturalsimulator
toperformtheirresearchstudies,themostcommonstepistomodifythesimulatortofittheneedsof
theirintendedstudy. Wefollowthe“Open-ClosedPrinciple”[80]—oneshouldbeabletoextenda
simulatorwithoutmodifyingit. Weavoiddesigninginter-dependentcomponents,becausemodifying
them tends to be non-trivial, and even more problematic than implementing them from scratch.
Instead,wedefineclearinterfacesandallowuserstowritenewextensionstoimplementtheirown
components. Theycanthenplugthenewextensionsintotheexistingsimulatorwithoutinterfering
withotherparts.
Accuracy and Trustworthy. Accuracy is one of the three primary goals during the
development of an architectural simulator [11]. While accuracy can be measured in comparison
toexistingsilicon,trustworthinessreflectstheabilityofthesimulatortomodelnovelideasnotyet
bakedintosilicon[28].
33

CHAPTER4. NAVISIM
4.2 Simulator Development Methodology
To achieve the design goals we mention above, we leverage the Akita simulator frame-
work[81]toenablethemodularandextensibledesignofvariousGPUcomponents. Akitaisbased
onadiscreteeventsimulation(DES)modelandconsistsofthefollowingkeyconcepts:
1. Theevent-enginesystem
2. Thecomponentsystem
3. Themessage-port-connectionsystem
Akita features an event-driven simulation engine which manages a queue of events for the entire
simulation and executes events sequentially. An event indicates a state update of the simulator.
Everyentityofthesimulatedcomputerplatformisacomponent[81]. TheMessage-Port-Connection
systemisusedforcommunicationbetweencomponents. Twocomponentscanonlycommunicate
witheachotherbysendingandreceivingmessagesthroughaconnection. Aconnectionisanabstract
structurethatcanmodelwires,buses,andintra-chip/inter-chipinterconnectnetworks.
TofulfillthegoalofdesigninganaccurateGPUsimulator,weusemicrobenchmarksto
calibratekeyparametersoftargetGPUsinordertoevaluatetheaccuracyofNaviSim. Wedesignaset
ofmicrobenchmarks(seeTable4.1)toevaluateindividualGPUsubsystems,suchasthewavefront
dispatcher,instructionpipelines,andmemoryhierarchy. Atahighlevel,werepeatoneoperation
thousands to millions of times to stress individual components of the GPU. Then we statistically
analyzethelatencyofanoperation,asweevaluatethethroughputofaspecificGPUcomponent.
Asanexampleofthisprocess,weusethepointer-chasingbenchmarktoevaluatecache
sizesandassociatedlatencies. Thehostprogramfirstcreatesaregionofmemory(i.e.,thepointer
chasingregion)ofagivensize. Next,thehostprogramdividestheregioninto8-bytecellsandfills
eachcellwithanaddressthatpointstothenextcell,locatedatarandomlyassignedaddresswithin
theregion. Thereisnorepetitionintheaddressesstoredintheregionsothatthewholeregioncan
be traversed multiple times. The GPU will use a single thread to access the cells many times (at
leastseveraltimesmorethanthenumberofcells),followingtheaddressesstoredintheregion. We
eventuallydividethekernelexecutiontimebythenumberofaccessestocalculatetheaverageaccess
latency.
Inthefollowingsectionsofthischapter,wewilldescribehowwemodelRDNAarchitecture
(ISAandmicroarchitecture)inNaviSimincludingtheimplementationdetailsandvalidationprocess.
34

| CHAPTER4. NAVISIM |     |     |     |     |     |     |     |     |
| ----------------- | --- | --- | --- | --- | --- | --- | --- | --- |
Table4.1: ThelistofmicrobenchmarksusedtocalibrateNaviSim.
| Microbenchmark |     | Parameter | Description |     |     |     |     |     |
| -------------- | --- | --------- | ----------- | --- | --- | --- | --- | --- |
empty kernel Work-Group Executesakernel,whilevaryingthekerneland
|     |     | DispatchingSpeed | work-groupsizes. |     | Weleavethekernelblankso |     |     |     |
| --- | --- | ---------------- | ---------------- | --- | ----------------------- | --- | --- | --- |
thatwecanmeasurework-groupdispatching
latency.
empty kernel multi KernelLaunching Launches a given number of empty kernels.
|     |     | Overhead | The kernel | used | is same | as  | the one | used in |
| --- | --- | -------- | ---------- | ---- | ------- | --- | ------- | ------- |
empty kernelmicrobenchmark.
single thread loop InstructionPipeline Executes a kernel with only one thread. The
|     |     | Depth | kernel | has a main | loop, | and | each iteration |     |
| --- | --- | ----- | ------ | ---------- | ----- | --- | -------------- | --- |
oftheloopexecutesasinglesingle-precision
floatingpointinstruction.
multi thread loop InstructionPipeline Executes the same kernel as
|     |     | Throughput | single | thread | loop, |     | but with | a   |
| --- | --- | ---------- | ------ | ------ | ----- | --- | -------- | --- |
largenumberofthreadsthatarelargeenough
tofullyoccupytheGPUexecutionresources.
pointer chasing random CacheSizesand Runstheclassicpointerchasingmicrobench-
|     |     | CacheLatencies | mark[82].      | Thepointersstoredinthebuffers |                        |     |     |     |
| --- | --- | -------------- | -------------- | ----------------------------- | ---------------------- | --- | --- | --- |
|     |     |                | arerandomized. |                               | Weonlyrun1threadinthis |     |     |     |
kernel.
| pointer chasing | linear | CacheSizeand   | The     | same    |     | kernel  |     | as     |
| --------------- | ------ | -------------- | ------- | ------- | --- | ------- | --- | ------ |
|                 |        |                | pointer | chasing |     | random, |     |        |
|                 |        | CacheLatencies |         |         |     |         |     | except |
thatthepointersinthebuffersalwayspointto
thenextpointer(thelastpointerpointstothe
firstone.
memory copy Bandwidthof Copiesdatafromonebuffertoanother,using
|     |     | CachesandDRAMs | agivennumberofwork-groups. |                  |          |                    | Whenthere |          |
| --- | --- | -------------- | -------------------------- | ---------------- | -------- | ------------------ | --------- | -------- |
|     |     |                | are only                   | 1-2 work-groups, |          | we                 | evaluate  | the      |
|     |     |                | L0cachebandwidth.          |                  |          | Whenthereareenough |           |          |
|     |     |                | work-groups                | to               | fill the | whole              | GPU,      | this mi- |
crobenchmarkcanalsotesttheDRAMband-
width.
35

CHAPTER4. NAVISIM
WestudythearchitecturaldifferencesbetweenGCN3andRDNAandbuildtheRDNAmodelfrom
scratchwhilereusingsomecomponentsinMGPUSim[20]. Showcasingthisprocessasaconcrete
exampleofthegeneraldesignmethodologyweproposedwillbenefitthesimulatorcommunity,not
onlyfromasimulationtoolwitharangeofease-of-usefeaturesbutalsothedesignphilosophyof
achievingahighlyadaptableandaccurateGPUsimulator. Wefocusonourflexibleemulatordesign
whichcanbepluggedintotimingsimulationmodular,makingitpossibletomixandmatchinstruction
sets with different microarchitecture implementations. This flexibility enables the evaluation the
performanceimpactofmicroarchitecturaldesignfeaturesintheRDNAarchitecture.
4.3 RDNA Architecture
TheAMDGCNarchitecture[2]hasbeenthereferencespecificationforAMDGPUdesign
overthepastdecade. However,inrecentyears,theAMDGCNarchitecturefacedcriticalscalability
challenges, as it struggled to fully exploit the ever-increasing number of transistors provided by
today’s photolithography technology [83]. In response, AMD developed a new family of RDNA
architectures[3],introducingamajorarchitecturaloverhaulascomparedtotheGCNarchitecture.
We focus on these architectural differences between GCN and RDNA architectures, including
the programming model, the compute unit (CU), and the memory hierarchy. The AMD RDNA
Whitepaper[3]providesmoredetailedspecificationsofthedesign.
OneofthemostnoteworthychangesisthattheRDNAarchitecturereducesthesizeofthe
wavefrontfrom64work-itemsto32work-items. Bycuttingthewavefrontsizeinhalf,theCUsare
expectedtobettercopewithahigherdegreeofthreaddivergenceinmodernworkloads. Additionally,
asthewavefrontsizeissmaller,fewermemorytransactionsareexpectedtobegeneratedbyasingle
load/storeinstruction(althoughthetotalnumberoftransactionsgeneratedbythekernelislikelyto
remainunchanged), potentiallyreducingmemoryaccesslatencies. Thesebenefitsshouldhelpto
improveALUutilization. Wewillfurtherevaluatetheeffectsofnarrowerwavefrontsinthelaterpart
ofthischapter.
AsecondmajorchangeinRDNAGPUsistheintroductionofDualComputeUnits(DCUs,
seeFigure4.1(d)),replacingtheGCNCUs. ADCUcontains4schedulers. Increasingthenumber
ofschedulersfrom1to4significantlyincreasestheinstructionissuerate. Ratherthandispatching
instructionsacross4SIMDunits, asintheintegratedGCNCUdesign, eachRDNAschedulerin
aDCUdispatchesinstructionsto1SIMDunit. OneSIMDunitinaDCUhas32single-precision
ALUs,doublingthenumberofALUsinaCU.Workingtogetherwiththenarrower32work-item
36

| CHAPTER4. | NAVISIM           |         |                  |     |     |     |     |              |                   |              |     |     |
| --------- | ----------------- | ------- | ---------------- | --- | --- | --- | --- | ------------ | ----------------- | ------------ | --- | --- |
|           | Command Processor |         |                  |     |     |     |     |              | Command Processor |              |     |     |
|           | S h ad er         |   Array | Sha d e r   A rr | a y |     |     |     | Shader Array |                   | Shader Array |     |     |
| L1I       |                   |         |                  | L1I |     |     |     |              |                   |              |     |     |
Instruct io n  M e m In s t r u c t io n Mem Dual Compute Unit Dual Compute Unit
| L1V | Compute Unit |     | Compute Unit |     |     |     |     |     |     |     |     |     |
| --- | ------------ | --- | ------------ | --- | --- | --- | --- | --- | --- | --- | --- | --- |
L1V
|     |              |     |              |     |     |     | ehcaC 1L L0S | L0V L0V | L0I | L0S | L0V L0V | L0I ehcaC 1L |
| --- | ------------ | --- | ------------ | --- | --- | --- | ------------ | ------- | --- | --- | ------- | ------------ |
| L1V | Compute Unit |     | Compute Unit | L1V |     |     |              |         |     |     |         |              |
L1V Compute Unit Compute Unit L1V Dual Compute Unit Dual Compute Unit
|     |              |     |              |     |     |     | L0S                | L0V L0V | L0I | L0S | L0V L0V            | L0I |
| --- | ------------ | --- | ------------ | --- | --- | --- | ------------------ | ------- | --- | --- | ------------------ | --- |
| L1V | Compute Unit |     | Compute Unit | L1V |     |     |                    |         |     |     |                    |     |
|     |              |     |              |     |     |     |            ... ... |         |     |     |            ... ... |     |
| L1S | Scalar Mem   |     | Scalar Mem   | L1S |     |     |                    |         |     |     | 3-5 DCUs in Total  |     |
3-5 DCUs in Total
| Interconnect |     |     |     | Interconnect |     | Interconnect |     |     |     |     | Interconnect |     |
| ------------ | --- | --- | --- | ------------ | --- | ------------ | --- | --- | --- | --- | ------------ | --- |
L2 Cache L2 Cache L2 Cache L2 Cache ... ... L2 Cache L2 Cache L2 Cache L2 Cache L2 Cache ... ... L2 Cache
DRAM DRAM DRAM DRAM ... ... DRAM DRAM DRAM DRAM DRAM ... ... DRAM
(a) The Architecture of GCN3  GPUs.
(b) The Architecture of RDNA GPUs.
| L1I      |                |     | L1S           | L1V        |      |      |            |           | L0I L0S |            |           |      |
| -------- | -------------- | --- | ------------- | ---------- | ---- | ---- | ---------- | --------- | ------- | ---------- | --------- | ---- |
|          | To Other CUs   |     |               |            |      |      |            |           |         |            |           | L0V  |
| 10 Wave  |                |     | Memory Global |            |      |      |            |           |         |            |           |      |
| Inst Buf |                |     | Access        | Register   |      |      |            |           | 20      |            |           |      |
|          |                |     |               | Write Back | SALU | SReg | VReg eussI | 20  hcteF | hcteF   | eussI VReg | SReg SALU |      |
|          | eussI & edoceD |     |               |            |      |      |            | Wave      | Wave    |            |           | Vec. |
hcteF tsnI 10 Wave Vec.  LDS Operand Gathering Inst Inst Operand Gathering Mem
| Inst Buf | VReg | SReg | SALU Mem | Unit |          |     |     | Buf | Buf |     |          |          |
| -------- | ---- | ---- | -------- | ---- | -------- | --- | --- | --- | --- | --- | -------- | -------- |
|          |      |      |          |      | SIMD-32  |     |     |     |     |     | SIMD-32  | tinU SDL |
Global
| 10 Wave  |     | Operand Gathering |     |     |        | Memory Access |     |     |     |     |          |     |
| -------- | --- | ----------------- | --- | --- | ------ | ------------- | --- | --- | --- | --- | -------- | --- |
| Inst Buf |     |                   |     |     |  D-32  |               |     |     |     |     | SIMD-32  |     |
|          |     |                   |     |     |        |               |     | 20  | 20  |     |          |     |
10 Wave SIMD S I M D S I M D S I M D Operand Gathering eussI Wave hcteF hcteF Wave eussI Operand Gathering Ve c .
| Inst Buf |     | -16  - 1 | 6   - 1 6 |   - 1 6   |      |      |      | Inst | Inst |      |           | M e m |
| -------- | --- | -------- | --------- | --------- | ---- | ---- | ---- | ---- | ---- | ---- | --------- | ----- |
|          |     |          |           |           | SALU | SReg | VReg |      | Buf  | VReg | SReg SALU |       |
Buf
L0V
(c) The Architecture of a GCN3 Compute Unit. (d) The Architecture of an RDNA Dual Compute Unit.
Figure4.1: AComparisonbetweentheGCNArchitecture[2]andtheRDNAArchitecture[3].
wavefronts,eachSIMDunitcanfinishexecutingoneinstructioninasinglecycle,ascomparedto4
cyclesinaCU.
Third, the RDNA architecture redefines the cache hierarchy from a 2-layer structure to
a3-layerstructure. ThecachesthataredirectlyconnectedtotheDCUsarerenamedasL0caches
(versusL1caches). Eachread-onlyL0instructioncacheandread-onlyL0scalarcachearenolonger
sharedbymultipleCUs,butdedicatedtoaDCU.EachDCUconnectswith2separatewrite-through
L0vectorcaches;agroupoftwoschedulersandtwoSIMDunitscanuseoneL0vectorcache. Since
nowwehavetwoL0cachesconnectedtooneDCU,updatingthedatainonecachemayrenderthe
dataintheotherL0cachestale. ThismaycausecoherenceissueswithinaDCUandrequiresexplicit
cacheinvalidationinstructions(asprovidedintheRDNAISA,butnotintheGCNISA).
Additionally,anintermediatelevelofcaching(i.e.,thenewwrite-evictL1cache)isinserted.
TheL1cacheservesagroupofDCUs(typically4-5)inaShaderArray,andsitsbetweentheL0and
L2caches. TheL1cachescanreducethenumberofrequestsarrivingattheL2caches(inthecaseof
L1hits)andreducetheamountofdatathatistransmittedacrossthechip(fromL2toL0),thereby
increasingperformanceandloweringthepowerconsumptioncausedbycross-chiptransmissions.
37

CHAPTER4. NAVISIM
Finally,thecachelinesizeoftheL0vectorcaches,L1caches,andtheL2cachesisdoubledfrom
64Bto128B,sothatacachelinecandeliveruniquesingle-precisionnumbersforall32work-items
inawavefront(4B×32=128B).
4.4 Implementation of NaviSim
Next, we focus on the implementation details including ISA emulation, wavefront dis-
patching,DCUmodeling,andmemoryhierarchy. Wealsodiscusstheuserinterface,thesupported
features,andsimulationoutputformats.
4.4.1 ISAEmulation
NaviSimrunsineitheremulationmodeortimingsimulationmode. Thecorecomponent
ininstructionemulationmodeisanemulatorwhichrecreatesexecutionresults,withoutevaluating
detailedtiminginformationfortheinstructionpipelines,caches,andDRAMcontrollers. Timing
simulationmodeutilizesanexternalemulatortoexecuteinstructionswhilesimulatingthetiming
behavior of other components. We have developed RDNA and CDNA instruction emulators for
NaviSim. Since these instruction emulators share the same interface with the GCN3 instruction
emulatorfromMGPUSim,thisallowsuserstoswaptheinstructionemulatorbeingused. Thevirtual
driver(asetofAPIsthatconnecttheGo-codedhostprogramsandthesimulatedGPUs)ofNaviSim
allowsuserstoconfigurewhichISAtoemulateandloadthecorrespondingGCN3/RDNA/CDNA
kernelbinaries.
Currently, NaviSim supports kernel written in either the OpenCL [10] or HIP [8, 9]
programming languages. OpenCL kernels can be compiled by the AMD official clang-ocl
compiler, which is a standard part of the AMD Radeon Open Compute (ROCm) platform [84].
HIP kernels can be compiled with the hipcc compiler, which also ships with the AMD ROCm
platform. Byusingthe--gencoargument,hipccignoresthehostprogramandonlygenerates
kernelbinaries. NaviSimsupportsloadingkernelbinariescompiledbyeithercompilerandemulates
theexecutionoftheGPUkernelsusingactualinputdatasets.
4.4.2 WavefrontDispatching
WemodelthewavefrontdispatchingprocessintheCommandProcessor. Themodeled
CommandProcessormaintainsresourcemasksthatkeeptrackofwhichresourcesareoccupiedin
38

CHAPTER4. NAVISIM
eachCU/DCU,includingthewavefrontslot(wavefront-levelresources,suchasinstructionbuffers
andtheprogramcounterregister),scalarregisters,vectorregisters,andLDSmemory. Theresource
masksensurethatnohardwareresourcesareoversubscribed. Sincewemasktheresourceseither
attheregisterlevel(forvectorandscalarregisters)oratthebytelevel(fortheLDS),wecanalso
modelregister/LDSfragmentationissuesintheCUs/DCUs[85].
Tosupportthemodelingofconcurrentkernelexecution[86],basicallyallowingmultiple
kernel executions simultaneously on one GPU, the modeled Command Processor has multiple
wavefrontdispatchers. Eachwavefrontdispatcherwillmanagetheprogressofthecurrentlyexecuting
kernel’s execution and dispatch a new wavefront when resources free up in the CUs/DCUs. By
default,eachCommandProcessorprovides8dispatchers(thenumberisconfigurable),andhence,we
allowupto8kernelstoexecuteconcurrently. Thewavefrontdispatcherscompetefortheresources
andhaveequalopportunitiestodispatchwavefrontstotheCUs/DCUs,ensuringthefairnessofthe
concurrentlyexecutingkernels.
4.4.3 DCUmodeling
WedevelopadetailedarchitecturalmodeloftheDCU,whichgovernshowinstructions
areexecutedinthesimulator. SincetheaccuratemodelingoftheDCUisessentialforsimulation
accuracy, NaviSim carefully models the pipeline (see Figure 4.2) with a multi-stage, multi-issue
structure.
Scheduler
retibrA
hcteF
... ...
20 Wavefront
Inst Buffers
in Total
... ...
retibrA
eussI
Decode Read Exec Write
Branch
Scalar
From 32-Width
L0 Inst ... ... ... SIMD
LDS
... Vector
Mem
50 Stages
To L0 Vector
Figure4.2: TheinstructionpipelinemodelinaNaviSimDCU.EachDCUhasfourcopiesofthe
structureshown.
Thepipelinestartswiththeinstructionfetcharbiter,whichattemptstofetchinstructions
39

CHAPTER4. NAVISIM
forawavefrontthathasvacantspaceininstructionbuffers. Whenmultiplewavefrontshavespace
available,thewavefrontthatreceivedinstructionsfurthestinthepastisselected. Theissuearbiter
monitors the instruction buffers and selects wavefronts that have instructions ready that they can
beissued. Thearbitercanissueatmost5instructionspercycle,oneinstructiontoeachinstruction
pipeline.
TheBranch,ScalarandLDSpipelinesuseafairlyregular6-stagepipelinethatincludes
fetch,issue,decode,read,execute,andwritestages. Notethat,forasingleinstruction,thedecode
stageislogicallyexecutedaftertheissuestage. Thisisbecausetheissuearbitercaneasilydetermine
thetypeofeachinstructionbycheckingafewencodedbitsintheinstruction. Thebranchinstruction
has no decode stage, as the instructions are very simple. This is in line with publicly available
documentation provided by AMD for the DCU architecture. The SIMD unit, which provides
most of the computing power of the DCU, also uses a 6-stage pipeline design, but is capable of
reading,executing,andwriting32instructionsinparallel. Thevectormemoryhasthemostcomplex
pipelinestructure,adoptinga50-stagepipeline. Thismaysoundunusual,butthemodelmatchesour
microbenchmarkingresults. Sincethereisnopublicinformationontheexactfunctionofeachofthe
pipelinestage,wedonotmodelthebehaviorofeachstage,butonlyassociatealatencyvaluewith
eachinstruction.
4.4.4 MemoryHierarchy
Wemodelthe3-levelcachehierarchyinRDNAGPUs. L0scalarandL0instructioncaches
areread-only,L0vectorcachesuseawrite-throughpolicy,L1cachesuseawrite-evictpolicy,and
L2cachesuseawrite-backpolicy. WealsoconnectDRAMSim3[87]tomodelGDDR5,GDDR5X,
GDDR6,GDDR6X,HBM,andHBM2DRAMcontrollers. Weallowuserstoconfigureanynumber
ofcachesinthehierarchyandsupportanycombinationoftheavailablecachepolicies. Additionally,
alloftheparametersofacache(e.g.,setcount,wayassociativity,cachelinesize,directorylatency,
storageaccesslatency)canbefullyconfiguredbyusers.
4.4.5 UserInterface
NaviSimadoptsasimilaruserinterfaceasusedinMGPUSim. Torunbenchmarks,theuser
cancompileeitherHIPorOpenCLkernelsusingtheofficialAMDcompiler. Ahostprogramwritten
inGoisrequiredtoinvokeasetofAPIstoallocate/copymemoryandlaunchkernels. ThehostAPIs
arecompatiblewithMGPUSimandaresimilartocommonGPUprogrammingframeworks(e.g.,
40

| CHAPTER4. | NAVISIM |     |
| --------- | ------- | --- |
Table4.2: Theconfigurationsofvalidationplatforms.
| Parameter    | Platform1      | Platform2      |
| ------------ | -------------- | -------------- |
| GPU          | RadeonRX5500XT | RadeonProW6800 |
| GPUCoreFreq  | 1845MHz        | 2320MHz        |
| DCUCount     | 11             | 30             |
| GPUMemory    | GDDR6          | GDDR6          |
| MemBandwidth | 224.0GB/s      | 512.0GB/s      |
AMDRyzen
| CPU |     | AMDEPYC7302P |
| --- | --- | ------------ |
Threadripper2950X
| OS        | LinuxUbuntu18.04 | LinuxUbuntu20.04 |
| --------- | ---------------- | ---------------- |
| GPUDriver | AMDROCm5.0       | AMDROCm5.1       |
CUDA,OpenCL,HIP).Userscanspecifyeitheremulationmode(valueemulationonly)ortiming
mode(detailedtimingsimulationthatinvolvesinstructionpipeline,caches,andmemorycontrollers)
asacommandlineoption. Meanwhile,NaviSimusesafewconfigurationfilesthatarewrittenin
Gotodefinethehardwareconfiguration. Userscaneasilyconfigurethehardwareundersimulation
bymodifyingparametersandcomponentconnectionsintheconfigurationcode. Weusecodefor
configurationbecauseuserscaneasilydebugtheconfigurationlogicwithdebuggers.
4.5 Evaluation
Next, we describe the methodology of validating NaviSim against real GPU hardware
platforms. Wecalibrateandvalidatetheaccuracyofthesimulatorusingasetofmicrobenchmarks
and full benchmarks. To demonstrate the model’s suitability to conduct a range of studies and
analyses,weconductadditionalexperimentstoevaluatetheRDNAarchitecturedesign.
4.5.1 SimulatorValidation
Anycomputerarchitecturesimulatorrequiresarigorousvalidationprocessbeforeitcan
serve as a baseline for future research. We present our validation work for NaviSim, comparing
simulationresultsagainstGPUhardwareexecution.
Validation Platforms: We use an RX 5500 XT GPU and Radeon Pro W6800 (see
Table4.2)tovalidateourNaviSimGPUModel. RX5500XTGPUisanentry-levelgraphicscard
41

| CHAPTER4. | NAVISIM |     |     |     |     |
| --------- | ------- | --- | --- | --- | --- |
whichwasannouncedalongwiththenewRDNAarchitecture. RadeonProW6800isahigh-end
workstationGPU,builtontheRDNA2architecture. ThemaindifferencesbetweentheRDNAand
theRDNA2lieinthememory/cachesystem(RadeonProW6800hasa4-levelcachehierarchy). Our
platformsruntheROCm5.0/5.1softwarestacksonLinuxUbuntu18.04/20.04servers.
Table4.3: Simulatorconfiguration
P3*
| Param      | P1       | P2      |     |         |     |
| ---------- | -------- | ------- | --- | ------- | --- |
| BaseModel  | RX5500XT | R9Nano  |     | -       |     |
| #CU/DCU    | 11(DCU)  | 64(CU)  |     | 32(DCU) |     |
| CoreFreq   | 1845MHz  | 1000MHz |     | 1000MHz |     |
| TFLOPS     | 5.20     | 8.19    |     | 8.19    |     |
| L0V$@      | 16KB     | 16KB    |     | 16KB    |     |
| L0V$Assoc. | 4-way    | 4-way   |     | 4-way   |     |
| L0Inst$@   | 32KB     | 32KB    |     | 32KB    |     |
L0Scalar$@
|              | 16KB    | 16KB       |     | 16KB    |     |
| ------------ | ------- | ---------- | --- | ------- | --- |
| L0I/L0S$Org. | PerDCU  | 4-CUshared |     | PerDCU  |     |
| L1$          | 128KB   | -          |     | 128KB   |     |
| L1$Assoc.    | 16-way  | -          |     | 16-way  |     |
| L2$          | 1MB     | 2MB        |     | 2MB     |     |
| L2$Assoc.    | 16-way  | 16-way     |     | 16-way  |     |
| DRAMTech     | GDDR6   | HBM        |     | HBM     |     |
| DRAMSize     | 4GB     | 4GB        |     | 4GB     |     |
| MemFreq      | 1750MHz | 500MHz     |     | 500MHz  |     |
| MemBus       | 128bit  | 4096bit    |     | 4096bit |     |
| MemBW        | 224GB/s | 512GB/s    |     | 512GB/s |     |
* P3isnotanoff-the-shelfGPU,butitischosenspecificallytoanalyzetheperformanceimpactof
individualRDNAfeatures.
@WewillrefertothecachesthataredirectlyconnectedtoaCU/DCUastheL0cache.
InP2,L0
cachesconnecttoL2cachesdirectly.
Simulator Configuration:
|     |     | We select | the baseline GPU | configuration | using publicly |
| --- | --- | --------- | ---------------- | ------------- | -------------- |
availableinformationandcalibrateourresultsusingmicrobenchmarks. Thedefaultconfigurationof
NaviSimfortheRX5500XTGPUsisrecordedasP1inTable4.3.
42

CHAPTER4. NAVISIM
AftervalidatingtheNaviSimGPUmodel,weusethemodeltoconductaseriesofuse-case
experimentstoevaluatetheperformanceimpactofmicroarchitecturaldesignfeaturesintheRDNA
architecture. We compare our simulation results with the default configuration of MGPUSim for
theR9NanoGPUs(denotedasP2inTable4.3). SincetheAMDRX5500XTandtheR9Nano
GPUsbelongtotwoverydifferentmarkets,comparingthemdirectlywouldnotprovideuswitha
lotofnewinsights. Therefore,wehaveconfiguredaGPUasP3inTable4.3. WechosetheseGPU
configurationssothattheywouldhaveacomparabletheoreticalcomputingthroughput(represented
byTFLOPS)andacomparablememoryhierarchy. Intheexperiments,ratherthandirectlycomparing
P2andP3,wegraduallyaddfeaturesfromP2tobuildP3,sothatwecanevaluatetheeffectofeach
feature.
FullBenchmarks: WealsoutilizealistoffullGPUbenchmarks(seeTable4.4)selected
fromawiderangeofbenchmarksuitesincludingAMDAPPSDK[88],SHOC[89],HeteroMark[90],
PolyBench[1],andDNNMark[91]. Weusethesebenchmarksasasetofworkloadsthatcoverawide
rangeofapplicationstoexercisethesimulatorwithdifferentarithmeticintensities,memoryaccess
patterns,andcommunicationpatterns. WestartwiththeoriginalOpenCLkernelimplementations
fromthebenchmarksuiteandcompilethekernelswiththeoriginalAMDROCmcompiler(applying
default compiler optimizations). We also write host programs in Go, allowing the simulators to
callthekernels. Weensureourhostprogramsareequivalenttotheoriginalhostprogramfromthe
benchmarksuite. Forvalidationexperiments,wevarytheproblemsizesofthebenchmarkstomake
surethatNaviSimcanrecreatescalingtrends. Forperformanceevaluation,weuselargeproblem
sizesthataresufficienttostressthewholeGPU.
Running Benchmarks: We run OpenCL implementations of the benchmarks, while
varying the input size. We use kernel execution time as the performance metric and report the
averageexecutiontimeover10runs. ThetimesobtainedontherealGPUarerecordedusingOpenCL
events[10].
WeusekernelexecutiontimeastheprimarymetrictoevaluatetheaccuracyofNaviSim,
for two reasons. First, execution time is the most commonly reported metric when considering
architecturaltradeoffs. Second,thekernelexecutiontimeisahigh-levelmetricthatsummarizesthe
impactofallfeaturesbeingsimulated;wecanonlyaccuracycaptureexecutiontimeifthefidelityof
allthecomponentsbeingmodeledinthesimulatorishigh.
Validation Steps: First, we verify the correctness (in terms of application outputs) of
NaviSim in both emulation mode and timing simulation modes. To this aim, we compare every
simulator-generatedapplicationoutputwithitscorrespondingactualhardwareexecutionoutput. We
43

CHAPTER4. NAVISIM
Table4.4: FullBenchmarks
Abbr. Suite Workload
ATAX PolyBench MatrixTransposeandVectorMultiplication
BICG PolyBench BiCGStabLinearSolver[92]
BS AMDAPPSDK BitonicSort
FIR HeteroMark FiniteImpulseResponseFilter
FLW AMDAPPSDK Floyd-WarshallAlgorithm
FWT AMDAPPSDK FastWalshTransform
KM HeteroMark KMeansClustering
MT AMDAPPSDK MatrixTranspose
ReLU DNNMark RectifiedLinearUnit
SPMV SHOC SparseMatrix-VectorMultiplication
observethattheresultsmatchexactlyforallapplications. ThematchingresultsofNaviSimwhen
runningintheemulationmodesuggestthatNaviSimfaithfullyemulatestheRDNAinstructions.
Second,tovalidatethetimingmodel,weusebothmicrobenchmarksandfullapplication
workloadstovalidatetheaccuracy. Weconfigureoursimulationusingthepubliclyavailableparame-
tersoftheAMDRX5500XTGPU(seeTable4.2). Wethenusemicrobenchmarks(seeTable4.1)to
helpusreverse-engineerawiderangeofundocumentedparameters,suchasthework-groupdispatch
overhead,instructionpipelinedepth,cachelatencyateachlevel,andDRAMbandwidthandlatency.
Agoodexampleofourstrategyisthatweuseourpointerchasingmicrobenchmarktofigureoutthe
sizeandlatencyofeachlevelofthecachehierarchy. Aswecanseefromtherepresentativeresults
ofrunningthepointerchasingmicrobenchmark(seeFigure4.3),NaviSimisabletocalibratethe
parameterswithextremelyhighaccuracy. Thecalibrationresultsofothermicrobenchmarksfollow
similartrendsasthepointerchasingresults.
Withthefullycalibratedparameters,weevaluatetheaccuracyofthesimulatorbyvalidating
itagainstAMDRX5500XThardwareusingfullbenchmarks(seeTable4.4). Figure4.4shows
thesimulatoraccuracy(lefty-axis),aswellastherelativeerror,ateachdatapointbythebarplot
(righty-axis)forarangeofproblemsizes. Theproblemsizerangesarespecifictoeachbenchmark.
Weobservethatforsomebenchmarkstherelativeerrorincreasesastheproblemsizeincreases(e.g.,
ReLU),andforsomeotherbenchmarkstherelativeerrordecreasesastheproblemsizeincreases
(e.g.,BS,SPMV).Foreachbenchmark,inparenthesis,wealsoreporttheaverageerror. Aswecan
44

CHAPTER4. NAVISIM
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
NaviSim RX 5500 XT GPU
Figure4.3: ComparisonofthePointerChasingmicrobenchmarkbetweenNaviSimandtheRX5500
XThardware.
see, theerrorintermsofmodeledexecutiontime, averagedacrossallbenchmarks, isjust9.75%.
Additionally,wenotonlymodeltheexecutionwithhighfidelity,weclearlycapturethepatternsand
thenuancesoftheGPUarchitecture. Forexample,intheFIRandReLUbenchmarks,wesuccessfully
capturethetransitionintheworkloadwhentheexecutiontimestartstoincrease. Thissuggeststhat
weproperlymodelanybottlenecksduetoresourcecontention(i.e.,numberofworkgroupsavailable)
ontheGPU.
Asanotherexample,weobservestep-wisebehaviorinsomebenchmarks(e.g.,BSbefore
32K, FWT after 16K and 32K, KM between 2K and 3K). As the observed step-wise behavior is
caused by complex interactions between the instruction scheduler, cache hierarchy, and memory
transactionhandling,beingabletomodelthesestepsdemonstratesthatNaviSimcanmodelthesubtle
featuresintheRDNAarchitecture. WehavealsovalidatedNaviSimagainstasecondGPUmodel,
theAMDRadeonProW6800GPU(seeTable4.2),whichisanRDNA2-basedGPU.Theaverage
differencebetweenthesimulatedandhardwaremeasuredexecutiontimesis10.08%. Ourvalidation
results,whichareshowninFigure4.5,demonstratethatNaviSimcancapturethemicroarchitectures
inRDNAacrossdifferentdeviceswithhighfidelity.
Wearealsoawareofdiscrepanciesinafewbenchmarks,suchasFWandKM.Ingeneral,
thesebenchmarksareeithershort-runningbenchmarks(FW)orworkloadsthatinvolvealargenumber
ofkernellaunches(KM).ThediscrepanciessuggestthatNaviSimhasdifficultiesinmodelingGPU
behaviorduringthekernellaunchphase. Ingeneral,webelievethisisnotabigproblemsincethe
simulatorsarelikelytobeusedtomodellargeproblemsizes,andthedifferencesobservedinthe
45

| CHAPTER4. NAVISIM |     |     |     |     |     |     |
| ----------------- | --- | --- | --- | --- | --- | --- |
kernellaunchoverheadshouldnotimpacttheoverallaccuracybymuch. Weleavemoredetailed
modelingofthekernellaunchbehaviorasfuturework.
MemoryFootprintandSimulationSpeed: Wealsoanalyzethememoryfootprintand
performance of NaviSim simulation. In the FIR benchmark (with a 4M problem size, as shown
inFigure4.4),weuse823MBofmemory,whichfitsinthememoryofmostmoderncomputers. In
termsofsimulationperformance,onanAppleM1MacMini,weachieve43.5KIPSand89.5KIPS
inserialandparallelmodes(NaviSimismulti-threaded),respectively. Thisperformanceismuch
faster than MGPUSim, which reported 27KIPS parallel execution performance in their original
paper. Thus,thememoryconsumptionandperformanceofNaviSimarequitecompetitive.
|             |     | Simulation  |     | Native Execution | Relative Error |                    |
| ----------- | --- | ----------- | --- | ---------------- | -------------- | ------------------ |
|             |     |             |     |                  | 60             | )%( rorrE evitaleR |
| ATAX (7.9%) |     | BICG (5.9%) |     |                  | 12 FIR (12.9%) | 16                 |
| 300         |     | 300         |     | 400              |                |                    |
| )su(emiT    | 12  |             |     | 16               | 9 40           |                    |
|             |     |             |     | 300              |                | 12                 |
| 200         | 9   | 200         |     | 12               |                |                    |
|             |     |             |     | 200              | 6              | 8                  |
|             | 6   |             |     | 8                | 20             |                    |
| 100         |     | 100         |     | 100 BS (3.9%)    | 3              | 4                  |
|             | 3   |             |     | 4                |                |                    |
| 0           | 0   | 0           |     | 0 0              | 0 0            | 0                  |
64 512 1K 1.5K 2K 64 512 1K 1.5K 2K 4K 16K 32K 48K 64K 64 1K 16K 256K
|              |     |     |     | 120 | 300           | )%( rorrE evitaleR |
| ------------ | --- | --- | --- | --- | ------------- | ------------------ |
| 8 FW (19.0%) |     |     |     | 8   | 16 MT (13.5%) |                    |
32 60
| )su(emiT 6 |     |               |     | 6 90            | 200 | 32  |
| ---------- | --- | ------------- | --- | --------------- | --- | --- |
|            | 24  |               |     |                 | 12  |     |
| 4          |     | 40            |     | 60              |     | 24  |
|            | 16  |               |     | 4               | 8   |     |
|            |     |               |     |                 | 100 | 16  |
| 2          | 8   | 20 FWT (5.5%) |     | 2 30 KM (11.5%) | 4   |     |
8
| 0   | 0   | 0   |     | 0 0 | 0 0 | 0   |
| --- | --- | --- | --- | --- | --- | --- |
16 64 128 192 256 4K 16K 32K 48K 64K 256 1K 2K 3K 4K 256 1K 1.5K 2K
| 60          | 16  | 100         |     | 40  |     |     |
| ----------- | --- | ----------- | --- | --- | --- | --- |
| ReLU (3.1%) |     | SPMV (9.8%) |     |     |     |     |
| )su(emiT    | 12  |             |     | 30  |     |     |
| 40          |     | 75          |     |     |     |     |
|             | 8   | 50          |     | 20  |     |     |
20
|                | 4   | 25     |     | 10   |     |     |
| -------------- | --- | ------ | --- | ---- | --- | --- |
| 0              | 0   | 0      |     | 0    |     |     |
| 64 512 8K 128K | 1M  | 64 512 | 1K  | 1.5K |     |     |
Figure4.4: SimulatorvalidationagainsttheAMDRX5500XTGPU.Thex-axisplotstheproblem
size and the two y-axes plot the kernel execution time and relative error. The numbers in the
parenthesesrepresenttheaveragedifferencebetweenNaviSimsimulationandthehardwareexecution
foreachapplication.
4.6 Case Studies: Understanding the RDNA Architecture Features
With a carefully validated simulator model, we next use NaviSim to perform a set of
experimentstoanalyzetheimpactofRDNAfeaturesonapplicationperformance. Inparticular,we
attempttoanswerthefollowingquestions:
46

| CHAPTER4. | NAVISIM     |     |              |     |                  |                |             |                    |
| --------- | ----------- | --- | ------------ | --- | ---------------- | -------------- | ----------- | ------------------ |
|           |             |     | Simulation   |     | Native Execution | Relative Error |             |                    |
|           | 300         |     |              |     | 1200             |                | 8 20        | )%( rorrE evitaleR |
|           |             |     | 300          |     |                  |                |             | 30                 |
|           | ATAX (9.3%) | 30  | BICG (11.4%) |     |                  |                | FIR (11.0%) |                    |
| )su(emiT  |             |     |              |     | 120              |                | 6 15        |                    |
|           | 200         |     |              |     | 800              |                |             |                    |
|           |             | 20  | 200          |     | 80               |                |             | 20                 |
4 10
|     | 100 |     |     |     | 400 |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
|     |     | 10  | 100 |     | 40  |     | 2 5 | 10  |
BS (3.3%)
|     | 0   | 0   | 0   |     | 0 0 |     | 0 0 | 0   |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
64 512 1K 1.5K 2K 64 512 1K 1.5K 2K 4K 16K 32K 48K 64K 64 1K 16K 256K
|     | 8        |     |     |     | 8 160 |     | 60            | 40 )%( rorrE evitaleR |
| --- | -------- | --- | --- | --- | ----- | --- | ------------- | --------------------- |
|     |          | 32  |     |     |       |     | 16 MT (14.6%) |                       |
|     | )su(emiT |     | 90  |     | 6 120 |     |               | 30                    |
|     |          | 24  |     |     |       |     | 12 40         |                       |
|     | 4        |     | 60  |     | 4 80  |     |               | 20                    |
|     |          | 16  |     |     |       |     | 8             |                       |
20
|     | FW (8.9%) |     | 30 FWT (4.2%) |     | 2 40 KM (7.0%) |     | 4   | 10  |
| --- | --------- | --- | ------------- | --- | -------------- | --- | --- | --- |
8
|     | 0   | 0   | 0   |     | 0 0 |     | 0 0 | 0   |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
16 64 128 192 256 4K 16K 32K 48K 64K 256 1K 2K 3K 4K 256 1K 1.5K 2K
|          | 30          | 15  | 300          |     | 100 |     |     |     |
| -------- | ----------- | --- | ------------ | --- | --- | --- | --- | --- |
|          | ReLU (6.3%) |     | SPMV (24.7%) |     |     |     |     |     |
| )su(emiT | 20          | 10  | 200          |     |     |     |     |     |
60
|     | 10  | 5   | 100 |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
20
0
|     | 0              | 0   |        |     | 0    |     |     |     |
| --- | -------------- | --- | ------ | --- | ---- | --- | --- | --- |
|     | 64 512 8K 128K | 1M  | 64 512 | 1K  | 1.5K |     |     |     |
Figure4.5: SimulatorvalidationagainsttheAMDRadeonProW6800GPU.Thex-axisplotsthe
problem size and the two y-axes plot the kernel execution time and relative error. The numbers
intheparenthesesrepresenttheaveragedifferencebetweenNaviSimsimulationandthehardware
executionforeachapplication.
1. HowdoeschangesintheISAimpacttheoverallperformanceandhowdoestheDCUarchitec-
tureimpactperformance? Howistheinstructionexecutionpipelineimpactedbyexecutinga
differentISA?
2. WhatistheeffectofthenewlyaddedL1cache?
3. Whatpercentoftheoverallperformanceincreasecanbeattributedtochangesinfrequency
(increasedfrom1GHzto1.845GHz)?
While we study the impact of the features on the overall performance, we also use this study to
demonstratetheutilityofNaviSim. WeshowcasetheflexibilityandtheconfigurabilityofNaviSim.
4.6.1 ImpactofchangesintheISAandtheintroductionoftheDCU
In the first set of experiments, we focus on question (1). We use the P2 configuration
(seeTable4.3,denotedasCU+GFX803inFigure4.6)asthebaseline. Wehaveconfigurationswhere
weeitherchangetheISAtoGFX1010(CU+GFX1010)orchangethecoretoDCU(DCU+GFX803,
where the number of DCUs is half of the number in the CU, thus providing a fair comparison of
47

CHAPTER4. NAVISIM
similarcomputingcapabilities). Wealsoprovideaconfigurationthatchangesboththecoreandthe
ISA(DCU+GFX1010).
1.5
1.0
0.5
ATAXBICG BS FIR FLW FWT KM MT PR ReLUSPMV
Benchmark
dezilamroN
emiT
noitucexE
CU+GFX803 DCU+GFX803
CU+GFX1010 DCU+GFX1010
Figure 4.6: The effect of changing from CUs to DCUs and the ISA from GFX803 (GCN3) to
GFX1010(RDNA)onexecutiontime.
Atahighlevel,ourresults(seeFigure4.6)suggestthatmanybenchmarks(e.g.,BS,PR,
ReLU)achievethesameperformanceforallfourcases. Thisisunderstandable,asthememorysystem
and the overall computing capabilities remain the same. However, we notice major performance
differencesinATAX,BICG,FIR,andFLWbenchmarks. Thesechangesarecausedbythedifferences
intheISAandtheCU/DCUorganization.
ATAX and BICG are good examples of how the CU/DCU organizer can impact the
performanceofworkloadsthataresensitivetomemorybandwidthandlatency. ATAXandBICGare
workloadswithlimitedparallelismandstronginter-work-itemdependencies. Therefore,ATAXand
BICGhavelargework-groupsizesandasmallnumberofwork-groupsineachkernel. IntheCU
configurations,thenumberofblockscannotutilizealltheCUsintheGPU,andhence,cannotfully
utilizethebandwidthbetweentheL0cachesandtheCUs. Onthecontrary,thenumberofblocksin
theATAXandBICGcanfullyutilizetheDCUsonthethirdandfourthconfigurations,sincethecount
ofDCUsishalved. AseachDCUisconnectedwithtwovectorL0caches,thebandwidthbetween
thecoreandtheL0cachesiseffectivelydoubled,resultingina2timesincreaseintheperformance.
ThemostcriticalchangeintheISAisthedifferenceinwavefrontsize. Kernelscompiled
to the GFX1010 ISA always use a wavefront size of 32, which is not a perfect match for the
CU architecture. When running wavefronts on CUs, because the CU scheduler can only issue 1
instructiontoaSIMDunitevery4cycles,theCUneedstospend4cyclestoexecute1instruction.
We observe underutilization of the ALUs in the ATAX, BICG, FIR, and FLW benchmarks when
48

CHAPTER4. NAVISIM
we run the GFX1010 ISA with the CU microarchitecture. Other benchmarks do not observe this
issuebecausetheyareboundbymemorybandwidthandarenotsensitivetoALUutilization. Inthe
FIRbenchmark,wealsoseethatrunningtheGFX803benchmarkontheDCUmicroarchitecture
causesasignificantslowdown. TheGFX803ISAusesawavefrontsizeof64andittakestwocycles
for the SIMD unit to execute, while the scheduler in the DCU schedule issue 1 instruction every
cycle. Thus,theSIMDunitisnotabletofinishexecutionoftheinstructionsintimetocatchupwith
theinstructiondispatch,leadingtoasignificantnumberofpipelinestallsandevenmoreslowdown.
TheproblemdoesnotexistwhentheDCUmicroarchitectureandGFX1010ISAareusedtogether,
suggestingthattheISAandthemicroarchitectureareco-designedtoachievethebestperformance.
4.6.2 ImpactoftheL1cache
Next,wefocusonquestion(2)andtrytounderstandtheeffectofthenewL1cache. We
use the P3 configuration (see Table 4.2, denoted as P3 w/ L1 in Figure 4.7) as the baseline, and
compare the performance with a configuration without an L1 cache (P3 w/o L1). The results are
showninFigure4.7.
1.0
0.5
0.0
ATAXBICG BS FIR FLW FWT KM MT PR ReLUSPMV
Benchmark
dezilamroN
emiT
noitucexE
P3 w/o L1 P3 w/ L1
Figure4.7: Theeffectofaddingthenewwrite-evictL1cacheonbenchmarkexecutiontime.
Similartotheearlierexperiments,asweonlymodifyasmallpartoftheconfiguration,most
ofthebenchmarksthatareeitherboundedbytheavailablecomputeresourcesorbyDRAMbandwidth
donotobserveaperformancedifference. However,wedidobserveperformanceimprovementsin
theFIRandFLWbenchmarks.
TounderstandwhyFIRandFLWcanbenefitfromtheaddedL1cache,weplotthecache
hitrateforeachlevelofcacheinFigure4.8. Ingeneral,weseeahighdegreeofdiversityacrossthe
benchmarks,astheseworkloadshavedramaticallydifferentlevelsoflocality. FortheBS,FIR,FLW
49

CHAPTER4. NAVISIM
andFWTbenchmarks,weobservearelativelyhighL1cachehitrate(around40%). However,the
effectofthehighL1cachehitrateontheL2cachesdiffersacrossworkloads. ForBSandFWT,the
L2cachehitrateremainsunchanged,whiletheL2cachehitratedropssignificantlyforFIRandFLW
(inagoodway). TheseresultssuggestthattheimprovedL1cachehitratesleadtoareductioninthe
numberofL2cachetransactionsfortheFIRandtheFLWbenchmarks. BecausetheFIRbenchmark
already has a high L0 hit rate, the speedup offered by the L1 cache is not as high as the benefits
observedintheFLWbenchmark.
XATA GCIB SB RIF WLF TWF MK TM RP ULeR VMPS
1.0
0.8
0.6
0.4
0.2
0.0
L0 Cache
etar
tiH
P3 w/o L1 P3 w/ L1
XATA GCIB SB RIF WLF TWF MK TM RP ULeR VMPS
1.0
0.8
0.6
0.4
0.2
0.0
L1 Cache
XATA GCIB SB RIF WLF TWF MK TM RP ULeR VMPS
1.0
0.8
0.6
0.4
0.2
0.0
L2 Cache
Figure4.8: Thecachehitrateofeachlevelofcache,beforeandaftertheL1cacheisadded.
4.6.3 Impactofincreasedfrequency
OnemajorchangeinRDNAGPUsistheincreasedfrequency. Whiletheincreaseisnotan
architecturalfeature,ahigherfrequencycanimpactotherarchitecturalchanges,soweareeagerto
understandhowmuchtheincreasedfrequencycanimproveperformance.
Here,weuseP3asabaselineandincreasetheGPUcorefrequencyto1845MHz,whichis
thebasefrequencyfortheRX5500XTGPU.Notethatasweincreasethecorefrequency,wealso
increasetheclockspeedthatcontrolstheL0,L1,andL2cachesto1845MHz. InAMDGPUs,the
coresandthecachesworkinthesamefrequencydomain. WedonotchangetheDRAMconfiguration,
sotheDRAMlatencyandbandwidthremainunchanged.
Overall,weseetheperformanceofallthebenchmarksimprove,withtherateofchange
rangingfrom1.25×to1.84×(seeFigure4.9). Thissuggeststhattheincreaseinfrequencyleads
toamarkedimprovementinGPUperformance. ThisisparticularlyevidentinthePRand ReLU
benchmarks, which experience a speedup which is closely correlated with the increase in the
core frequency. This suggests that the kernel execution time of these two benchmarks is linearly
50

CHAPTER4. NAVISIM
2.0
1.5
1.0
0.0
ATAXBICG BS FIR FLW FWT KM MT PR ReLUSPMV
Benchmark
pudeepS
P3@1000MHz P3@1845MHz
1.5
Geo-Mean
Figure4.9: Theeffectoftheincreasedcorefrequencyontheoverallperformance. Asweincrease
thecorefrequencyby1.845×,theperformanceimprovedby≈1.5×.
dependent on the GPU frequency. More specifically, after deeper investigation with our visual
simulationdebugger(Daisen[93],wefindthatthetwoworkloadsaredependentonthework-group
dispatchrate,andhavearathershortwork-groupexecutiontime. Thedispatchercannotcatchupwith
thework-groupretirementrate. BecausetheCUcanretirethework-groupsfasterthantheratethatat
whichnewwork-groupsarrivetostartexecution,thenumberofwork-groupsconcurrentlyexecuting
in each CU is limited. This also limits the number of concurrent in-flight memory transactions.
Increasing the core frequency accelerates the dispatching speed and increases the occupancy of
the CUs, allowing the CUs to better utilize the memory bandwidth. This effect is also observed
in the case study discussed in the Daisen paper [93]. Meanwhile, we also notice that the average
speedup(geometricalmean)isonly1.5×,suggestingthattheincreasedcorefrequencystillneedsto
beaccompaniedbyimprovedDRAMperformance.
51

| Chapter | 5     |                    |     |              |
| ------- | ----- | ------------------ | --- | ------------ |
| Design  | Space | Exploration:       |     | Accelerating |
| FHE     | with  | Microarchitectural |     | Extensions   |
(GME)
Withafullyvalidatedsimulatormodel,wenextuseNaviSimtocharacterizeapplication
behaviorandexplorenoveldesignideasusingNaviSim. Wewilltargettheperformanceissueofthe
FHEworkloadonGPUs. WewilldemonstratetheflexibilityofNaviSim,attemptingtoacceleratea
FullyHomomorphicEncryption(FHE)workload.
Our approach to accelerating the workloads starts by profiling the application on our
simulator. Onceweunderstandwhereperformancebottlenecksoccur,wefocusonenhancingthe
microarchitectureofGPUsandwewillexplorethreedifferentmicroarchitectureextensions,aswell
asacompile-timeoptimization,toaccelerateFHEexecutiononanAMDCDNAGPU.
To evaluate our design, we extend the implementation of NaviSim with a block-level
directedacycliccomputegraphsimulatorcalledBlockSim[69]. Inaddition, weconductablation
studies on our micro-architectural feature implementations, enabling us to isolate each micro-
architecturalcomponentandevaluateitsdistinctinfluenceontheentireFHEworkload.
| 5.1 Challenges |     | of Accelerating | CKKS FHE |     |
| -------------- | --- | --------------- | -------- | --- |
TheCKKSscheme,similartootherhomomorphicencryptionschemes,introducescompu-
tationaloverheadduetoitsuniqueproperties. Inordertoperformcomputationsonencrypteddata,
52

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
largevaluesofN andQareusedforencoding/decodingandencryption/decryption,whichproduce
heavyworkloaddemands. Itrequiresasignificantamountofcomputationalresourcestoevaluate
modulararithmeticexpressionsandalsoplaceshighdemandsonmemorybandwidthutilization. We
firstconductaworkloadcharacterizationoftheGPUkernelsthatimplementmultipleFHEprimitives.
As an example of this process, we use the number-theoretic transform (NTT) kernels, which are
essential building blocks of polynomial multiplication, to reveal the intrinsic bottlenecks of FHE
workloads. Figure5.1providesprofilingresultsfortheNTTkernelsforintegersizesof32bitsand
64bits. Thehighmemoryintensityreportedindicatesthatasignificantportionoftheexecutiontime
isdevotedtoaccessingmemory,resultinginlowALUutilization. Weobservethat,notonlyNTT,
butalsoalloftheFHEoperationsexhibitlowarithmeticintensityandrequireverylargeworking
sets,whichindicatesfrequentmainmemoryaccesses.
|     |     |     | 32-bit |     | 64-bit |     |     |
| --- | --- | --- | ------ | --- | ------ | --- | --- |
100
)%( egatnecreP
60
20
0
|     | n   |     | e   |     | e   |     | nsity |
| --- | --- | --- | --- | --- | --- | --- | ----- |
|     | o   |     | Rat |     | Rat |     |       |
ati
|     | utiliz |     | Hit  |     | Hit  |     | nt e  |
| --- | ------ | --- | ---- | --- | ---- | --- | ----- |
|     |        |     | e    |     | e    |     | ory I |
| U   |        |     | h    |     | h    |     |       |
| AL  |        | ac  |      | ac  |      |     |       |
|     |        | C   |      | C   |      | m   |       |
|     |        | 1   |      | 2   |      | M e |       |
|     |        | L   |      | L   |      |     |       |
Figure5.1: NTTkernelsprofileofALUutilization,cachehitratesandmemoryintensity
Memory Bottlenecks: The basic datatype in CKKS is ciphertext, denoted as a vector
N logQ
of elements, where each element is bits in size. Thus, we can calculate the size of each
ciphertext. We choose N = 216 and logQ = 1728, meaning that our ciphertext size will be 2∗
N∗logQ=28.3 MB, where each polynomial in the ciphertext is ∼14 MB. After applying RNS
53

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
decomposition[45]onthesepolynomialsusingawordlengthof54bits,weget32limbsineach
polynomial,whereeachlimbis∼0.44MBlarge.
ForAMD’sMI100GPUs,thesizeofthelastlevelcacheandlocaldatasharespaceare
8MBand7.5MB,respectively. Thuswecannotaccommodateevenasingleciphertextintheon-chip
memory. Atmost,wecanfit∼18limbsofaciphertextpolynomial,andasaresult,wewillhaveto
performfrequentaccessestothemainmemorytooperateonasingleciphertext.
Compute Bottlenecks: FHE operations involve complex mathematical computations,
includingmodulararithmeticandpolynomialmultiplications. Modulararithmeticoperationsoccur
atthelowestlevel,whichmeanseveryscalaradditionandmultiplicationisfollowedbyamodular
reduction. Sincedivisionoperationisinvolvedandisnotdirectlyimplementedonhardware,some
algorithms have to be used to perform a modular reduction operation, so it becomes one of the
compute bottlenecks. In addition, we need to use RNS decomposition to map the wide modular
operations to the widths of the underlying computing platform (i.e., 32, 64, and 128-bit) and
thenperformoperationsonthedividedsmallerelements(RNSlimbs). Thelimbbitwidthimpact
the number of modular addition and multiplication operations required. Currently, commercially
availableGPUarchitecturesdonotimplement64-bitand128-bitwidemodulararithmeticoperations
directlybutcanemulatethemviamultiplearithmeticinstructions,whichsignificantlyincreasesthe
amountofcomputingrequiredfortheseoperations.
To address the previously identified bottlenecks, we focus on exploring new micro-
architecturalextensionswithNavisimthatcanleveragetheGPU’son-chipmemoryspacetominimize
theaccessestothemainmemoryandprovidenativemodulararithmeticunitstoacceleratemodular
arithmeticoperations.
5.2 Exploring Microarchitectural Extensions with NaviSim
In this work, we propose to re-architect the current GPU microarchitecture and also
introducenovelmicroarchitecturalextensionsthatenableimprovedutilizationofGPUresourcesso
astomaximizetheperformanceoftheFHEworkloadsrunningontheGPU.Weproposearobustset
ofmicro-architecturalfeaturestargetingAMD’sCDNAarchitecture,unlockingthefullpotentialof
theGPUtoaccelerateFHEworkloads.
WepinpointcriticalbottlenecksencounteredduringFHEworkloadexecutionandaddress
themprogressivelyusingthreemicro-architecturalfeatureextensions. Ouron-chipCU-sidehierar-
chicalnetwork(cNoc)andtheLocalityAwareBlockScheduler(LABS)contributetominimizing
54

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
the DRAM bandwidth bottleneck. Simultaneously, our implementation of native modular reduc-
tion (MOD) and wider multiply-accumulate units (WMAC) features improve the math pipeline
throughput,ensuringastreamlineddataflowwithevenlydistributedresourceutilization.
5.2.1 cNoC:CU-sideinterconnect
ModernGPUshaveanetwork-on-chipthatinterconnectsthecomputeunits(CUs)together
with the memory partitions or memory banks (see Figure 2.1). On-chip communication occurs
betweentheCUsandthememorybanks,notnecessarilybetweentheCUs.
V-ALU
Interconnect
Memory
Bank
L2 Cache
DRAM
1L ehcaC
Memory
Bank
L2 Cache
DRAM
SDL 1L ehcaC
Interconnect
Memory
Bank
L2 Cache
DRAM
1L ehcaC
Memory
Bank
L2 Cache
DRAM
SDL SDL SDL 1L ehcaC
V-ALU V-ALU V-ALU
R R
(a) Traditional Mem-hierarchy (b) On-chip routers allow data sharing
data sharing requires memory bypassing the off-chip interconnect
transactions to traverse the entire stack
Figure5.2: Inter-CUcommunication: Traditional(left)versusproposedon-chipnetwork.
In current GPU microarchitecture design, we can utilize the Local Data Share (LDS)
memoryspacededicatedtoonecomputeunittoenablelow-latencycommunicationbetweenthreads
inonework-group. However,thememoryaccesspatternoftheFHEworkloadrequiresdatatobe
sharedwithalargegroupofthreadswhosesizeismuchlargerthanthesizeofawork-group. Thus,
itisnotpossibleinthecurrentGPUmicroarchitecturedesigntosupportdatasharingbetweenwork-
groupsthroughCU-sidememoryspace(i.e.,LDS).WeproposalanovelCU-sidenetwork-on-chip
55

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
cNocthatinterconnectstheCUstogether,adesignthatisinspiredbytheideaoftheLDS.Weenable
a “global” LDS that can be shared between the CUs. We show how data is communicated in the
traditionalmethodandintheproposednewon-chipnetwork(seeFigure5.2)
ByexploitingthecNoc,thededicatedon-chipmemorycanbesharedbetweencompute
units, thus reducing the number of redundant memory accesses to the off-chip memory. We also
providesynchronizationbarriersofvaryinggranularitytomitigateraceconditionsbymodifyingthe
logicofthekerneldispatcher. SincetheLDSuseisunderprogrammercontrol,ourapproachdoes
notincurtheoverheadassociatedwithcachecoherenceandavoidsredundantcacheinvalidations,
butcomeswithsomeextraprogrammereffort. Byimplementingaglobaladdressspace(GAS)inour
GPU,weestablishdatasharingandformaunifiedGASbycombiningallLDSs. Thevirtualaddress
spaceisthenmappedontothisunifiedGAS,withtranslationusingahashoftheloweraddressbits.
Shader Array
... ...
Shader Array
Shader Array Shader Array
Shader Array Shader Engine
Shader Array Shader Array
rabssorC
pihC-ffO
Compute Shader Off-chip On-chip 2D Torus
LEGEND MUX
Unit Array Crossbar Router Topology
Figure5.3: Proposedhierarchicalon-chipnetworkfeaturingaconcentrated2Dtorustopology.
CurrentGPUsaredesignedhierarchically(e.g.,MI100GPUcomprisesnumerouscompute
units),with8CUscombinedtoformaShaderArray. Thereare15ShaderArraysintotal,resulting
in a total of 120compute units. Theproposed cNoc takes advantageof this hierarchy, utilizinga
hierarchical on-chip network (see Figure 5.3) that features a single router for each Shader Array,
connectingtheeightcomputeunitsthatmakeupaShaderArray. Theroutersarearrangedina3×5
2Dgridandinterconnectedthroughatorustopology. Whilethisconcentrated-torustopology[94,95]
canincreasenetworkcomplexity,itreducesthenumberofrequiredrouters(from120to15),thereby
minimizingthechipareaneededforthenetwork. Inaconcentrated-torustopology,allroutershave
56

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
thesamedegree(numberofports),creatinganedge-symmetrictopologythatiswell-suitedforthe
all-to-allcommunicationpatternsofFHEworkload.
Theon-chipnetworksignificantlyenhancesdatareuse,reducingredundantmemorytrans-
actions by 38%, and effectively supports the all-to-all communication pattern common in FHE
workloads.
Weimplementtheon-chipnetwork,leveragingtheMessage-Port-Connectionsystemin
NaviSim[81]. Weextendtheshaderarraystructurewithanon-chiprouter,leveragingsimplerouting
algorithms. We use a connection to connect CUs in the shader array with the router as a node
in the 2D torus network. We reuse the communication protocols of accessing the LDS unit for
read/writeoperations. Theprocessofimplementingthenetworkdoesnotrequiremodifyingfilesof
existingcomponents;onlyfileswithnewextensionsareadded,aswellaschangesinconfiguration
filesforpluggingthenewextensionsintotheexistingsimulatortorealizenewconfigurations. To
furtherexperimentandvalidatethedesign,NaviSimallowsuserstoeasilyconfigurethehardware
undersimulationbymodifyingparametersandcomponentconnectionsintheconfigurationcode.
NaviSim exhibits high adaptability in terms of effectively supporting extensions of the current
microarchitectureimplementations.
5.2.2 EnhancingtheVectorALU
Native modular reduction extension: The existing GPU arithmetic pipeline is highly
optimizedfordatamanipulationoperations,suchasmultiply,add,bit-shift,andcompare. Awavefront
executinganyoftheseinstructionstakes4clockcyclestoexecuteintheSIMDunitsinalock-step
fashion. Conversely, operationssuchasdivideandmodulusareemulatedusingaseriesofnative
instructions,resultinginpoorerperformancecomparedtotheirnativecounterparts.
Asstatedbefore,themodularreductionoperation,usedfordeterminingtheremainderofa
division,isperformedaftereachadditionormultiplication. Asaresult,optimizingmodularreduction
iscrucialforspeedingupFHEworkloads. Atpresent,theMI100GPUexecutesamodularoperation
throughasequenceofaddition,multiplication,bitshift,andconditionaloperations,drawingonthe
conventionalBarrett’sreductionalgorithm. Thisprocessisconsiderablyslow.
TodemonstratetheflexibilityofourNaviSimsimulationframework,wewillextendthe
vectorALUpipelineoftheCDNAarchitecturetonativelysupportmodularreduction. Weaugment
the CDNA instruction set architecture (ISA) with a collection of vector instructions designed to
performmodularreductionoperationsnatively,afterexecutinganadditionormultiplicationoperation.
57

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
Thenewnativemodularinstructionsproposedinclude:
• Nativemodularreduction:
mod-red <v0,s0> |V =V mods
0 0 0
• Nativemodularaddition:
mod-add <v0,v1,s0> |V =(V +V )mods
0 0 1 0
• Nativemodularmultiplication:
mod-mult <v0,v1,s0> |V =(V ×V )mods
0 0 1 0
Modularreductioninvolvesseveralcomparisonoperations,resultinginbranchdivergence
inGPUs. OurimplementationisderivedfromanimprovedBarrettreductionalgorithm[16]. This
approachminimizesthenumberofcomparisonoperationstoonepermodularreductionoperation,
significantlyreducingthenumberofbranchinstructionsandenhancingcomputeunitutilization.
Wider multiply-accumulate units In the CKKS FHE scheme [13], we can choose to
performoperationson32,64,or128-bitwideRNSlimbsforaciphertext. Thislimbbitwidthgoverns
theoperandsizeforthevectorALUs,impactingthenumberofmodularadditionandmultiplication
operationsrequired.
We are faced with an algorithm-level performance versus precision decision when se-
lectingthebitwidth. Ifweoptfor32-bitwideRNSlimbs,wewillhavenumerouslimbstowork
with,increasingtheavailablelevels,whilesimultaneouslyreducingtheachievableprecisionofthe
application. Conversely, if we select 128-bit RNS limbs, we will have fewer limbs to work with,
resultinginadecreaseinthenumberofavailablelevels,butinhighprecisionforanapplication. With
ourchosenparameters,using128-bitwideRNSlimbswouldleaveuswithaninsufficientnumber
of limbs to perform a single bootstrapping operation, which is a sequence of the building block
operationsTable2.1thatreducesthenoisewithintheciphertexttoanacceptablelevel. Tostrikea
balancebetweenperformanceandprecision,wechoosetouse64-bitwideRNSlimbsinthiswork.
MostGPUsonthemarketnativelysupport16-,32-,and64-bitfloatingpointcomputations,
aswellas4-,8-,32-bitintegercomputations. Unfortunately,theylackdedicatedhardwaresupport
for64-bitintegeroperations,themostcommonlyusedformatinFHEworkloads. Instructionsfor
processing64-bitintegeroperandsareemulatedusingmultiple32-bitintegerinstructions,making
themsignificantlyslower.
Tocomplementournativemodularreduction,whichrelieson64-bitintegeroperations,
weaddsupportforahardware-backed64-bitintegermultiplierandaccumulator. Wealsowidenthe
58

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
Table5.1: Cyclecountsfor64-bitmodulusinstructionscomparingMODandWMACfeatures
µ-arch. mod-red mod-add mod-mul
Feature (cycles)∗ (cycles)∗ (cycles)∗
VanillaMI100† 46 62 63
MOD∆ 26 18 38
MOD+WMAC 17 7 23
†ReferstotheunmodifiedCDNAarchitectureofMI100GPUs.
∗Cyclecountisaveragedover10,000modulusinstructions,computedoncacheddata(usingLDScache)androundedtothenearest
integer.
∆Modularoperationisimplementedusingvariouscompile-timeprimeconstantsasthemodulus.
register-filesizetoaccommodatethelargeciphertexts. Table5.1reportsthedecreaseintotalcycles
foreachofourproposednativemodularinstructions,incomparisontotheMI100GPU-emulated
instructionsinthebaselineconfiguration.
Modularreductioninvolvesexecutingseveralcomparisonoperations,resultinginbranch
divergence in GPUs. Our implementation is derived from an improved Barrett reduction algo-
rithm [16]. This approach minimizes the number of comparison operations to one per modular
reductionoperation,significantlyreducingthenumberofbranchinstructionsandenhancingcompute
unitutilization.
5.2.3 LABS:Locality-AwareBlockScheduler
Sofar,ourmicro-architecturalextensionsprimarilyfocusedonoptimizingindividualFHE
blocks. Tobetterleveragethesenewfeatures,wefocusnextoninter-blockoptimizationopportunities,
targeting the work-group dispatcher within the CDNA architecture. GPU scheduling is typically
managedusingstreamsofblocksthatarescheduledoncomputeunitsinagreedymanner[8]. The
presenceoflargeGPUregisterfilesallowstheschedulertooversubscribeblockstoeachcompute
unit. However,theexistingschedulerwithintheCDNAarchitectureisnotcognizantofinter-block
datadependencies,forcingcacheflusheswhentransitioningfromoneblocktothenext.
WeproposeaLocality-AwareBlockScheduler(LABS)designedtoscheduleFHEblocks
withshareddatatogether,thusavoidingredundanton-chipcacheflushes,specificallyintheLDS.The
LABSfurthercomplementsoursetofmicro-architecturalenhancements,whichrelaxtheoperational
constraints during block scheduling and create new opportunities for optimization. For example,
thecNocfeatureenablesLDSdatatobegloballyaccessibleacrossallCUs,therebyallowingthe
59

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
schedulertoassignFHEblockstoanyavailableCU.
TodevelopLABS,weemployawell-knowngraph-basedmappingsolutionandframethe
problemofblockmappingtoCUsasacompile-timeGraphPartitioningProblem(GPP)[96].
GraphPartitioningProblem: WeformulateourGraphPartitioningProblem(GPP)by
introducing a cost function Φ. For a graph G, if it is partitioned such that E denotes the set of
c
edgecuts,thenΦcanbeexpressedas|E |,thesumoftheindividualcut-edgeweights. Here(v,w)
c
representstheedge-weightoftheedgeconnectingnodevtonodew.
Φ=|E |= ∑ |(v,w)|
c
(v,w)∈Ec
Inthisequation,|(v,w)|signifiesthedatatransferredbetweenFHEblocks. Φisthecostfunction
thatreflectsthecommunicationoverheadassociatedwithassigningFHEblockstoseparatecompute
units. Thegoalofthegraphpartitioningproblemistodiscoverapartitionthatevenlydistributesthe
loadacrosseachcomputeunit,whileminimizingthecommunicationcostΦ.
Topartitionthecomputegraphandprepareitformappingontothearchitecturegraph,we
utilizeamultilevelmeshpartitioningtechnique. Forreadersinterestedingainingfurtherinsightsinto
ourgraphpartitioningimplementationofthemulti-levelmeshpartitioningalgorithm,werecommend
referringtotheworkofWalshawandCross[96].
Architecture-awaremapping: Weusetwographsinourlocality-awareblockscheduler.
LetG=G(V,E)representadirectedacycliccomputegraphwithverticesV andedgesE. HereV
correspondstoFHEblocksandE indicatesthedatadependenciesoftheFHEblocks. Similarly,let
G =G (V ,E )denoteanundirectedgraph,withverticesV andedgesE . Here,V representsGPU
a a a a a a a
computeunits,andE illustratesthecommunicationlinksbetweencomputeunits. Bothedgesets,E
a
andE ,areassumedtobeweighted,withedgeweightsofE signifyingthesizeofdatatransferred
a
betweenrelatedblocks,andE representingthelatencyofcommunicationbetweencorresponding
a
computeunits. Wecanthendefineπ :V →V asamappingofV intoV disjointsubsets.
a a
Our objective is to find a mapping π that minimizes communication overhead between
computeunits. WefocusonmappingourpartitionedsubgraphsontothesetofcomputeunitsV ,
a
wherecommunicationcosts(bothlatencyandbandwidth)arenotuniformlydistributedacrossthe
network. Touniformlydistributethecommunicationoverheadsacrossthenetwork,weintroducea
networkcostfunctionΓ. Formally,Γisdescribedas:
60

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
Γ= ∑ |(v,w)|.|(π(v),π(w))|
(v,w)∈Ec
Here,Γisdefinedastheproductofindividualcut-weightsandtheircorrespondingedge-
weightsinthearchitecturegraph,whenmappedusingamappingfunctionπ. Inthisequation,π(v)
represents the mapping of block v to a compute unit from the setV , after applying the mapping
a
function π. Additionally, |(π(v),π(w))| represents the communication latency between compute
units π(v) and π(w). Similar to our analysis with Φ, our goal is to minimize Γ. To accomplish
this, we use a compile-time optimization by applying simulated annealing [96], alongside mesh
partitioningforefficientmappingofFHEblocksontocomputeunits.
5.3 Evaluation
5.3.1 Methodology
We evaluate the three microarchitecture extensions using NaviSim, demonstrating its
flexibilityandextensibilitytosupportarchitecturaldesignspaceexploration. Toevaluatetheoverall
designincludingtheinter-FHE-blockoptimization,ablock-leveldirectedacycliccomputegraph
simulator called BlockSim (developed by Kaustubh Shivdikar) is introduced on top of NaviSim.
BlockSim[69]isdesignedtoenableustoidentifyinter-kerneloptimizationopportunities. Withan
adjustable sampling rate for performance metrics, BlockSim accelerates simulations, facilitating
moreefficientdesignspaceexploration. BlockSimgeneratesanalyticalmodelsoftheFHEBlocksto
provideestimatesforruntimesofvariousGPUconfigurations. Whenthebestdesignparameters
areidentified,NaviSimisthenemployedtogeneratecycle-accurateperformancemetrics. NaviSim
is highly configurable using a few configuration files that are written in Go to define hardware
configuration. AsagoodexampleofNaviSim’sdegreeofreconfigurability,BlockSimenablesthe
restructuringofthekerneldispatcherbymodifyingthedispatchinglogictointegratecompile-time
optimizations obtained from LABS. Besides supporting FHE workloads, BlockSim serves as an
essential component of NaviSim by abstracting low-level implementation details from the user,
allowing them to focus on entire workloads rather than individual kernels. We utilize AMD’s
CDNAarchitecture-basedMI100GPU(seetable5.2)tosetthebaselineGPUconfigurationforFHE
applicationevaluations.
We iteratively introduce micro-architectural extensions and evaluate the performance
benefits of each enhancement. We first evaluate our three micro-architectural extensions (cNoC,
61

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
Table5.2: MI100GPUParameters
Parameter Value
GPUCoreFreq 1502MHz
ProcessSize 7nm
ComputeUnits 120
TFLOPS 23.07
RegisterFile 15MB
CUcount 120
L1VectorCache 16KBperCU
L1ScalarCache 16KB
L1InstCache 32KB
SharedL2 8MB
LDS 7.5MB
GPUMemory 32GBHBM2
MemBandwidth 1229GB/s
HostCPU AMDEPYC7002
HostOS Ubuntu18.04
GPUDriver AMDROCm5.2.5
62

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
Table5.3: PerformanceofvariousFHEbuildingblocks
tluMC
ddA-EH
tluM-EH
etatoR
elacseR
BaselineMI100(µs) 178 217 4012 3473 681
GME∗ (µs) 22 28 464 364 69
SpeedupoverBaseline 8.1× 7.8× 8.6× 9.5× 9.9×
∗ThevaluesdisplayedhereexcludecontributionsfromtheLABSoptimization,asLABSisaninter-blockoptimization,andthemetrics
providedareintendedforindividualblocks.
MOD,WMAC),thenevaluateourcompile-timeoptimizationLABS,andconcludewithamemory
sizeexplorationtodeterminetheimpactofon-chipmemorysizeonFHEworkloads. Weevaluate
thesemicro-architecturalenhancementsandcompileroptimizationusingNaviSimandBlockSim.
Weanalyzetheperformanceofthreeworkloads: bootstrapping[43],HE-basedlogistic
regression(HELR)[97]andencryptedResNet-20[98]utilizingtheCIFAR-10dataset. Forallthree
workloads, we evaluate the contributions of each individual FHE building block (see Table 2.1)
thatmakeuptherespectiveworkload. Inaddition,fortheseworkloads,wereporttheperformance
benefitsachievedbyemployingeachoftheproposedmicro-architecturalenhancements.
5.3.2 Results
PerformanceofFHEBuildingBlocks: Webeginbyshowingtheperformanceofindi-
vidualFHEblocks. SincetheseareindividualFHEblocks,thereportedmetricsdonotaccountfor
ourinter-blockLABScompileroptimization. WefindthatHEMultandHERotatearethemost
expensiveoperations,astheyrequirekeyswitchingoperationsthatinvolvethemostdatatransfers
fromthemainmemory. ThenextmostexpensiveoperationisHERescale,wheretheruntimeis
dominatedbythecompute-intensiveNTToperations.
Across the five FHE blocks mentioned in Table 5.3, we achieve an average speedup
of 8.8× compared to the MI100 baseline execution. In particular, we see a substantial perfor-
mance improvement in the most expensive operations, namely HEMult and HERotate, as our
proposedmicro-architecturalenhancementsreducethedatatransfertimeby12×forbothblocks.
ForHERescale,wemanagetodecreasetheaveragememorytransactionlatencyby13×usingour
micro-architecturalenhancementstotheon-chipnetwork,cNoC.
Impact of Micro-architectural Extensions: Figures 5.4 and 5.5 highlight the impact
63

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
arch Feature
Bootstrapping
HE-LR
ResNet
Figure5.4: Influenceofindividualproposedmicro-architecturalextensiononarchitecturalperfor-
mance metrics. Metrics illustrate a cumulative profile where each enhancement builds upon the
precedingsetofimprovements
64

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
enilesaB CoNc DOM SBAL SDLx2
3.5
3.25
3
2.75
2.5
2.25
2
1.75
1.5
1.25
1
pudeepS
Bootstrap
enilesaB CoNc DOM SBAL SDLx2
HE-LR
enilesaB CoNc DOM SBAL SDLx2
ResNet
Figure 5.5: Speedup achieved from each micro-architectural extension. The baseline refers to a
vanillaMI100GPU.Thereportedspeedupiscumulative,witheachmicro-architecturalenhancement
buildinguponthepreviousones
ofeachofourproposedmicro-architecturalextensions,aswellasourcompile-timeoptimizations,
acrossthreedifferentworkloads,i)bootstrapping,ii)HE-LR,andiii)ResNet-20.
First, our proposed concentrated 2D torus network enables ciphertexts to be preserved
in on-chip memory across kernels, leading to a significant increase in compute unit utilization
acrossworkloads,therebyreducingtheaveragecyclesconsumedpermemorytransaction(seeAvg.
CPT in Figure 5.4). In fact, when comparing the average number of cycles spent per memory
transaction(averageCPT),weobservethattheResNet-20workloadconsistentlyachievesalower
average CPT value compared to the HE-LR workload. This indicates a higher degree of data
reuse within the ResNet-20 workload across FHE blocks, as opposed to the HE-LR workload.
WiththecNoCenhancement,asthedatarequiredfrompreviouskernelsisretainedintheon-chip
memory,CUsarenolongerstarvedfordataandthisalsoresultsinasubstantialdecreaseinDRAM
bandwidthutilizationandDRAMtraffic(thetotalamountofdatatransferredfromDRAM).The
L1cacheutilizationdecreasesnotablyacrossallthreeworkloadsforthecNoCmicro-architectural
enhancement. ThisisduetothefactthattheLDSbypassestheL1cache,andmemoryaccessesto
theLDSarenotincludedintheperformancemetricsoftheL1cache.
TheproposedMODextensionenhancestheCDNAISAbyaddingnewinstructions. These
newinstructionsarecomplexinstructionsthatimplementcommonlyusedoperationsinFHE,like
mod-red,mod-add,andmod-mult. Astheseinstructionsarecomplex(composedofmultiple
sub-instructions),theyconsumeahighernumberofcyclesascomparedtosimplerinstructionssuch
65

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
as mult or add. This is the reason for the increase in the average cycles per instruction (CPI)
reportedinFigure5.4.
Thecompile-timeLABSoptimizationinourapproachfurtherremovesredundantmemory
transactions by scheduling blocks that share data together, thus reducing total DRAM traffic and
enhancingCUutilization. LABStakesadvantageofon-chipciphertextpreservationfeatureenabled
byourcNoCmicro-architecturalenhancement. AsFigure5.5shows,acrossbootstrapping,HE-LR,
andResNet-20workloads,LABSconsistentlydeliversanadditionalspeedupofover1.5×ontopof
cNoCandMOD.
On-chipMemorySizeExploration: Finally,weexploretheimpactofon-chipmemory
(LDS)sizeonFHEworkloadefficiency,asshowninFigure5.6. ByincreasingthetotalLDSsize
from 7.5MB (which is the current LDS size on MI100 GPU) to 15.5MB, we achieve speedups
of 1.74×, 1.53×, and 1.51× for Bootstrapping, HE-LR, and ResNet-20 workloads, respectively.
However, increasingtheLDSsizebeyond15.5MBdoesnotresultinsubstantialspeedup, asthe
DRAMbandwidthbecomesabottleneck.
66

CHAPTER5. DESIGNSPACEEXPLORATION:ACCELERATINGFHEWITHMICROARCHITECTURALEXTENSIONS(GME)
GME
g
n
pi
p
a
r
t
s
t
o
o
B
GME
R
L
-
E
H
GME
0
2
-
t
e
N
s
e
R
Figure5.6: Exploringtheimpactofon-chipmemorysizeonFHEworkloadperformance.
67

Chapter 6
Simulation Performance Evaluation and
Optimization
So far in this dissertation we have focused on designing and building highly adaptable
simulators. Inthischapter,wehaveconsideredsimulationperformanceasthemaindesigncriteria.
Whenreferringtothetermperformancethroughouttheremainderofthischapter,wearereferringto
thespeedofthesimulation.
Cycle-levelarchitecturalsimulatorstendtobeslowduetotheirnatureofcapturingthe
detailedtimingbehavioroftheunderlyinghardware. Therequiredlevelofdetailcomesatthecostof
speed[22]. Multi2Sim[58]isreportedtobe44,000×slowerthannativeexecutiontime. Malhotra
etal.[61]reporttheslowdownoverthenativeexecutiontimeofGPGPUSimis480,000×. Based
onourevaluation,MGPUSimexperiencesaslowdownof5to6ordersofmagnitudecomparedto
nativeexecutiontime. TherapidevolutionofGPUarchitectures,whichhasbeendrivenbyHPCand
MLworkloads,requiresGPUarchitectstorapidlyevaluatetheimpactofavarietyofnewfeatures
withperformantsimulators. Thischallengedrivesthemotivationfordesigninghigh-performance
simulators.
Theexcessivelyslowsimulationspeedmotivatesustoreconsideroverallsimulatordesign,
withafocusonperformance. Ourgoalistodiscoverpatternsinthedesignprocessthatcontributeto
simulationoverhead. Someofthekeyfactorsincludethedegreeofabstractionselectedformodeling,
thesoftwareengineeringapproachemployed,theefficiencyoftheparallelsimulationengine,and
theimpactofapplyingcertainsimulationstrategies(i.e.,DES)haveonthesimulationperformance.
Todesignahigh-performancesimulator,understandinghowthesefactorsinfluencethesimulation
68

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
speedandexploringpossiblestrategiestomitigatesomeoverheadareessentialsteps.
6.1 Methodology
Next, we cover our methodology for examining how different aspects of the simulator
designaffectsimulationspeed. Designingahigh-performancesimulatorisanextremelychallenging
task because it involves multiple complex steps, including: selecting a simulation framework,
decidingonthebestabstractionlevelforthesimulatormodels,andadoptingasimulatorarchitecture
thatwilldrivetheimplementationofthemodel. Todesignahigh-performancesimulator,itmakes
sensetoexploitmulti-threadingtoacceleratearchitecturalsimulation,especiallygiventhenumberof
coresavailableontoday’sCPUs. Beingabletorunasimulationinparallelbecomesakeyfactorin
achievingfastandaccuratesimulationresults.
Given all these considerations, we narrow our discussion to focus on the key design
decisionsthatsimulatordevelopersneedtoconsidertominimizeunnecessarysimulationoverhead
andimproveparallelsimulationefficiency. Ourgoalistodevelopbestpracticesthatcanhelpguide
futuresimulatordevelopers.
6.1.1 Generalsimulatordesign
Sincediscreteeventsimulation[33](DES)isoneofthemostcommonlyandwidelyused
methodsofsimulationinthefieldofcomputerarchitecture,ourdiscussionfocusesonthedesignof
event-basedsimulators. Discreteeventsimulationmodelsthesystemasaseriesofeventsoccurring
atspecificmomentsintime. Thebehaviorofthesystemismodeledthroughhandlingevents.
Retry mechanism: Scheduling events to perform work in the future is a flexible and
convenient way to model the system behavior. While this approach to simulation has a number
of advantages, it still may encounter significant delays. For example, scheduling events that are
needed to retry a blocked operation is a common mechanism, and managing this through events
worksefficientlyinmostcases. However,insituationswheretherateofrequestsexceedstherate
ofrequestconsumption,especiallyaswescaleinputdatasetsizes,thisretrymechanismcancause
significant overhead. Since message congestion cannot be resolved immediately, continuing to
generateadditionaleventstoattemptthesameoperationimpactsimulatorperformanceandfailto
makeanyforwardsimulationprogress.
69

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
Toavoidthispitfall,weneedtoconsiderhowtohandlepotentialcongestionresultingfrom
alargevolumeofrequests. Weneedtodesignaspecificmechanismtoaddressitattheconsumer
end. Predictingallpossiblescenariosthatcouldleadtocongestioncanbechallenging. Toaddress
thisissueeffectively,weneedtodesignoursimulatortomonitorthetotalnumberofactiveevents
andidentifyeventsthatresultinmultipleretryoperations.
Emulation: Whendesigningsimulatorsthatsupportbothfunctional(i.e.,foremulation)
and timing simulation, a key design consideration is whether to create a standalone emulation
model or to develop a unified model that serves both emulation and timing simulation. While a
customemulationmodelmaybemoreperformant,developingasinglemodelthatcanusedforboth
emulationandtimingsimulationwillsupportdevelopingreusablecodeandincreasemaintainability
ofthesimulatorcode. However,itcouldintroduceperformancedrawbacksduringemulation.
For event-driven simulation, the granularity of events defined by the simulator greatly
impactsthesimulationspeed. Comparedtotimingsimulation,whichgenerallyrequiresmorefine-
grainedeventdefinitions,emulationcanleverageacoarsergrainedeventdefinitionthatcanboost
performance. Fromaperformanceperspective,aseparateemulationmodeloffersadvantagesovera
unifiedmodel.
Level of abstraction: When designing a simulator, a critical decision is to decide the
level of the detail we should use when simulating the hardware. This decision is critical since it
canimpactboththesimulationspeedandtheaccuracy. Thegoalofthearchitecturalsimulatoristo
enableperformancemodelingofthehardwarebehavior. Athoroughunderstandingofthefactors
thatimpacttheunderlyinghardwareperformanceiscrucialfordeterminingtheappropriatelevel
ofabstractionselectedforthesimulator. Developersmustnavigatethetrade-offbetweensimulator
performanceandaccuracytooptimizetheirdesigndecisions.
Anecessarystepistoexaminetheperformanceafterthesimulatorisbuiltandvalidated.
Wewillundoubtedlyfindthatsomeofourdecisionsleadtonotableperformancebottlenecks. Can
weraiseourlevelofabstractionbystrippingawaycertaindetailsinourmodelwhilestillmaintaining
accuracy? The answer depends on whether the details being removed significantly impact the
performance. Oneapproachistodevelopmultiplemodels,eachpossessingadifferentlevelofdetail
andperformexperimentstodecidewhichoneworksbest,intermsofaccuracyandperformance.
70

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
6.1.2 Parallelsimulationscalability
Moderncomputersystemsareinherentlyparallel. Giventhegrowingnumberofcoresin
currentCPUs,weshouldexplorehowbesttoutilizethesecorestorunsimulationefficiently. This
mayenableustosimulatetheexecutionofourdesignsinparallel. Althoughtheparallelsimulation
mechanismcanimproveperformanceascomparedwithserialsimulation,theperformancescalability
issuebecomesaconcern.
Parallelsimulationengine: Aparallelsimulationengineallowsmultiplesimulationtasks
torunconcurrently,distributingtheworkloadacrossmultipleprocessorsorcores. Forevent-driven
simulators,parallelismcanbeintroducedbyprocessingmultipleeventsconcurrently. Intheparallel
engine,weaggressivelymanagetheeventschedulingandexploitmultiplethreads,butwestillneed
to synchronize shared data access and ensure the correct order of operations. The efficiency of
the parallel simulation engine is crucial for achieving high scalability when performing parallel
simulation.
Synchronization: Whenimplementingcomponentsaspartofamulti-threadedprogram,
weneedtocarefullyexaminethedifferentscenarioswheresharedresourcesaremodifiedandproperly
managetheresources(withlocksorsemaphores)topreventdataraces. However,synchronization
overhead can significantly degrade performance. Finding a balance between an efficient locking
mechanism and scalability is the key to designing a parallel simulation framework. We believe
havingahighlymodulardesignforeachcomponentcankeepthelockingmechanismsimpleand
clean. Wealsobelieveseparatingtheinitializationofthesharedresourcesfromaccessingthemcan
reducethesynchronizationrequirements.
Connections: The design of the communication system is essential for the scalability
of parallel simulation. If the messages exchanged between components cannot be transferred in
parallel,thecommunicationsystembecomesamajorbottleneckforallthreads. Astraightforward
connectiondesignishavingauniversalcommunicationsystemconnectingallcomponents,which
could cause significant synchronization overhead because different messages interfere with each
other by sharing the same connection. A more parallel-friendly communication system should
haveindividualconnectionslinkingtwoports(ownedbycomponentsorswitches)andusemultiple
switchestomanagerouting.
71

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
6.1.3 PerformanceMetrics
Themostcommonmetrictomeasuresimulationspeedisoverallsimulationexecutiontime
(i.e.,wall-clocktime). Tobettermeasurethesimulationoverhead,weintroduceasetofperformance
metrics,thatofferacomprehensiveviewofthemeasurementofthesimulationoverhead.
Slowdown: Slowdownindicateshowmanytimesslowerthesimulationisascomparedtonative
execution. Thisisthemostintuitivemethodofmeasuringthesimulationoverhead.
simulationexecutiontime
Slowdown=
nativeexecutiontime
KIPS: Kilo-simulated-instructions per second (KIPS) indicates how many seconds are spent to
simulatetheexecutionof1000instructions. Thesimulatedinstructioncountisanimportant
metric that architectural simulators provide, in addition to the predicted kernel execution
time[58].
simulated instructioncounts
KIPS=
simulationexecutiontime∗1000
CIPGI: CPUinstructionspersimulatedGPUinstruction(CIPGI)indicateshowmanyinstructions
areexecutedonthehostsystemtosimulateoneGPUinstruction.
CPU instructioncounts
CIPSHI=
simulated GPU instructioncounts
6.2 Evaluation
To explore the various factors that influence simulator performance, and develop an
engineered process for tuning the performance of the simulator, we use NaviSim to ground our
discussion. WeconductathoroughperformanceanalysisofNaviSimwhilerunningsimulationsin
both emulation mode and timing simulation mode. We profile NaviSim as an application to find
bottlenecksandperformanceissues. Wediscussourfindingsandattempttopinpointbottlenecksand
inefficienciesintheexecutionofthesimulatorcode. Ourgoalistobuildaprocessthatleverages
the execution profile of NaviSim to help guide a set of software-level optimizations to improve
simulationspeed.
6.2.1 ExperimentalSetupandProfilingMethods
We evaluate a set of GPU benchmarks (see Table 6.1) selected from a wide range of
benchmarksuitesandsupportedbyNaviSim. WeelecttousethedefaultconfigurationoftheR9
72

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
NanoGPUasoursimulatorconfiguration(detailsareprovidedinTable6.2). Ourexperimentsrun
in two different simulation modes: i) emulation mode and ii) timing simulation mode, on Linux
platforms(seeTable6.3).
Table6.1: Benchmarks
Abbr. Suite Workload
ATAX PolyBench MatrixTransposeandVectorMultiplication
BS AMDAPPSDK BitonicSort
FWT AMDAPPSDK FastWalshTransform
FIR HeteroMark FiniteImpulseResponseFilter
FLW AMDAPPSDK Floyd-WarshallAlgorithm
KM HeteroMark KMeansClustering
MM AMDAPPSDK MatrixMultiplication
MT AMDAPPSDK MatrixTranspose
PR HeteroMark PageRank
ReLU DNNMark RectifiedLinearUnit
SPMV SHOC SparseMatrix-VectorMultiplication
To profile NaviSim’s execution, we utilize the Go performance profiling tools, which
includepprofandtrace[39],toperformathoroughanalysisasweidentifyperformancebottlenecks.
We use the Go profiling tool pprof to collect profiling data and visualize the different profiles,
includingheapusage,CPUusage,threadblockingprofilesandexecutiontracedata. TheCPUprofile
identifies the code sections that dominate overall execution time. The CPU profile provides both
a textual and graphics-based visualization of the profile data to help identify the most expensive
portions of the simulation. The Go profiling tool trace provides execution details of each Go’s
lightweightthread(calledgoroutine)bytracingtheruntimeeventswithinaninterval. Thegoroutine
analysisviewdisplaysinformationabouteachsetofgoroutinesthatsharethesamefunction. Since
NaviSimisanevent-drivensimulator,wealsoexploreandanalyzetheeventinformation(i.e.,event
name,eventtime,andexecutiontime)ofalltheeventsbeingexecutedinthesimulation,inspecting
them in time order. This helps us identify efficiencies during event handling and parallel event
execution.
73

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
Table6.2: Simulatorconfiguration
| Param      | Large   |     |
| ---------- | ------- | --- |
| BaseModel  | R9Nano  |     |
| #CU        | 64      |     |
| CoreFreq   | 1000MHz |     |
| TFLOPS     | 8.19    |     |
| L1V        | 16KB    |     |
| L0V$Assoc. | 4-way   |     |
L1Inst$@
32KB
| L1Scalar$@ | 16KB    |     |
| ---------- | ------- | --- |
| L2$        | 2MB     |     |
| L2$Assoc.  | 16-way  |     |
| DRAMTech   | HBM     |     |
| DRAMSize   | 4GB     |     |
| MemFreq    | 500MHz  |     |
| MemBus     | 4096bit |     |
| MemBW      | 512GB/s |     |
Table6.3: Evaluationplatform
| Parameter | Platform1 | Platform2   |
| --------- | --------- | ----------- |
|           | AMDRyzen  | AMDEPYC7742 |
CPU
|         | Threadripper2950X | 64-CoreProcessor |
| ------- | ----------------- | ---------------- |
| OS      | LinuxUbuntu18.04  | CentOS7Linux     |
| Cores   | 16C/32T           | 64C/128T         |
| MaxFreq | 3.5GHz            | 2.25GHz          |
| L1      | 32KB(D)+64KB(I)   | 32KB(D)+32KB(I)  |
| L2      | 512KB             | 512K             |
| L3      | 8MB               | 16MB             |
| DRAM    | 64GB              | 1TB              |
74

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
6.2.2 PerformanceAnalysis
Inthissection,wetakeadeeperlookattheperformanceofNaviSim. Thisanalysishelpsto
groundourmethodology,whichwillbeusedtomakekeydesigndecisionsthatsimulatordevelopers
needtoconsiderinordertotunesimulationperformanceandimproveparallelsimulationefficiency.
6.2.2.1 OverallPerformanceanalysis
Aswementionedinourmethodology,examiningtheperformanceafterthesimulatoris
builtandvalidatedisanecessarysteptooptimizethedesigndecisions. Withtheassistanceofthe
textualandgraphics-basedvisualizationoftheprofilingresults,weidentifythatakeybottleneck
involvesthebuilt-incopyfunction,whichinternallycopiesvaluesfromonebuffertoanotherby
calling the runtime.memmove function in Go’s runtime. Figure 6.1 is a partial CPU profile
generated from running the Go pprof tool on top of running the FIR benchmark. We find both
emulationandtimingsimulationencounterthisbottleneck.
Weareabletotracetheoriginofthisbehaviortolocationsinthesimulatorwherecopyfunc-
tionsarecalled,discoveringthataretiedtoaccessesinthescratchpadandscratchpadpreparer
datastructures. Thescratchpadisadatastructure,definedintheimplementationoftheALU
componentusedfortheALUtostoreinputandoutputdata. Thescratchpadpreparerisadata
structurethatpreparestheoperandsbyreadingfromorwritingintotheregisterfilesbeforeALU
startstoaccesstheseoperands. Thescratchpadisaninternalstoragespaceinthesimulatorto
modelsomehardwarelogicusedfortemporarilystoringoperandsfortheALU.
This performance bottleneck poses a question of whether we could raise the level of
abstractionwhenwemodeltheALU.Weoffersomepotentialsolutionsinsection6.3.
6.2.2.2 ParallelSimulationScalability
NaviSim supports simulating hardware operations in parallel by leveraging the Akita
simulatorframework[81]. Akitaemploysaparallelsimulationenginetosupportparallelsimulation.
The parallel engine triggers events that are scheduled at the same time in parallel. This strategy
followsaconservativeparalleldiscrete-eventsimulationmethodology(PDES)[99],whichmaintains
thechronologicalorderoftheeventsandwillnotimpactsimulationcorrectness.
WeevaluatetheefficiencyoftheAkitaparallelsimulationengineusedinNaviSim. First,
werunasetofbenchmarks,whileleveragingthecoresonourmulticoreCPUusedforsimulation.
75

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
Figure 6.1: The CPU profile from the Go pprof tool, while running the FIR benchmark timing
simulation.
76

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
WesetthenumberofCPUcorestobe1, 4,8, and16, whichiscontrolledusingtheenvironment
variableGOMAXPROCS.
Emulation: Figure 6.2 shows the speedup achieved through parallel simulation when
runningemulationmode. Wenormalizethesimulationexecutiontimeofrunningwith4,8,and16
coresagainst1coretoshowthespeedup. Theaveragespeedupsfor8and16coresare2.64and2.44.
Formostbenchmarks,runningwith8coresgivesthebestperformance,andwhenincreasingto16
cores,weeitherobtainthesamespeeduporperformancedegrades.
Clearly, there are some performance issues with the parallel emulation. Our approach
toaddresstheseperformanceissuesistostudytheeventlistduringemulation,withoutapplyinga
parallelengine. First,weanalyzetheproportionofeventsthatdotheactualinstructionemulation,
becausethesearetheeventsthatcanbehandledinparallel. Wefindthatonlyaportionofthetotal
time is spent in executing events that are related to emulation. The rest of the events are created
and handled by simulating other parts of the system (i.e. driver, dispatcher). We believe that the
WorkgroupDispatcheristhemostinefficientcomponent,andthereareopportunitiestostreamlineits
functionalitysoastoreducesomeunnecessaryevents.
6
4
2
0
ATAX BS FWT FIR FLW KM MM MT PR ReLUSPMV
Benchmark
pudeepS
4 8 16
2.65
2.44
2.2
Geo
Mean
Figure6.2: Thespeedupachievedthroughparallelsimulationwhenrunningemulationmode.
Apart from analyzing the event list, we use the Go tool trace to further investigate the
behaviorofthegoroutines. Figure6.3showsatableofspecificgoroutineinstances,withvarious
executionstatistics. Figure6.4showsaportionofasynchronizationblockingprofile. Theprofiling
results shown in Figure 6.3 and Figure 6.4 are generated when running the ReLU benchmark
emulationinparallelwithGOMAXPROCS=16. Thetimeintervalwelookathereis100ms. Thefull
77

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
synchronizationblockingprofilecanbefoundinFigureA.1.
The goroutine analysis view indicates only half of the total profiling time is spent on
execution for all goroutines, and about 20% of the time is spent as Block time caused by the
sync package. The synchronization blocking profile further confirms that these goroutines are
competing for shared resources (i.e., they are accessing the same data structure simultaneously).
The sync package provides locks to avoid race conditions. As a result, these requests must be
executed sequentially, which introduces waiting/blocking time. The graph of the profile data,
shown in Figure 6.4, visualizes where this blocking occurs. Each box in the graph corresponds
to a single function, and the boxes are sized according to the length of the time in which the
function spends running. The largest box is the Lock function from the sync package. We
see that synchronization blocking is implemented in three functions — find, getTable and
createOrGetStorageUnit. createOrGetStorageUnit is the function that spends the
mosttimewaiting.
Figure6.3: Executiontimebreakdowningoroutineanalysisview,usingtheGotracetool.
Timingsimulation: Figure6.5showsthespeedupachievedthroughparallelsimulation
when running in timing simulation mode. Similar to the emulation results, we normalize the
simulation execution time when running with 4, 8, and 16 cores against the time for a single
core, reporting simulation speedup. The average speedups for 8 and 16 cores are 2.12 and 1.97,
respectively. Weareencounteringasimilarissueaswefoundintheemulationresults. However,the
underlyingcauseoftheperformancescalabilityissueforsimulationisnotthesameastheemulation.
78

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
Figure6.4: ThesynchronizationblockingprofilefromtheGotracetool,whilerunningtheReLU
benchmarkparallelemulation.
79

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
2
1
0
ATAX BS FWT FIR FLW KM MM MT PR ReLUSPMV
Benchmark
pudeepS
4 8 16
2.12
1.97
1.63
Geo
Mean
Figure6.5: Thespeedupachievedthroughparallelsimulationwhenrunningtimingsimulation.
Comparedtoemulation,thenumberofeventshandledintimingsimulationgrowsexponentially. The
memorysubsystemsimulationgeneratesasignificantnumberofeventsduringtimingsimulation,
whichdiffersincontrasttoemulationmode,whichmodelsanidealmemorysinceemulationisonly
concernedwithcorrectness.
Figure6.6isaportionofthesynchronizationblockingprofilegeneratedwhenrunningthe
ReLUbenchmarktimingsimulationinparallelwithGOMAXPROCS=16. Thefullsynchronization
blockingprofilecanbefoundinthe FigureA.2.
OneofthekeyinterfacesofthesimulationengineisSchedule. Scheduleregistersevents
byinsertingthemintoaneventqueue. Runisanotherkeyinterfactthatextractseventsfromtheevent
queueandletstheeventsbehandled. Tosupportconcurrenteventhandling,weprovideseveralevent
queues in the parallel simulation engine. The engine manages them through Channel which is
abuilt-indatatypeinGoforconcurrencyandisusedforcommunicationbetweengoroutines. To
avoidraceconditions,modifyinganeventqueueonbothends(i.e.,insertingintothequeue(initiated
bySchedule),andextractingfromthequeue(signaledbyRun)simultaneously,isnotallowed. A
Channel is used to control when the event queues can be modified, enabling the goroutine to
scheduleevents,whilekeepingtheavailableeventqueuesinthechannel. Theprofileindicatesthat
thegoroutinesarewaitingtoreceiveaneventqueuefromthechannelqueueChanintheSchedule
function,whiletryingtoinsertaneventintothequeuetocompleteschedulingthatevent. Weobserve
that the first element in the queue array always receives more events than the other queue array
entries. Asthepositiontoenterthequeuebecomesdeeperinthequeuearray,thefewereventscan
80

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
Figure6.6: ThesynchronizationblockingprofilefromGotracetoolforReLUbenchmarkparallel
timingsimulation
81

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
be received. This is because the scheduling policy always tries to schedule events in the earliest
queuearrayentry. Thelargenumberofeventsassignedtothefirstqueuearrayentryimpactsthewait
timeforthenextroundofexecutions,sinceallgoroutineshavetowaitbeforethisnumberofevents
are launched, in order to have access to the first available queue. To lower the number of events
scheduledinthefirstqueuearrayentryineachround,weutilizeastrategybyswappingtheposition
oftheentriesinthequeuearray. Thelastqueuearrayentryassignedinthepreviousroundbecomes
the firstqueue array entry inthe current round. Thiscan reduce the queueprocessing time at the
beginningofeachround.
600
400
200
0
ATAX BS FWT FIR FLW KM MM MT PR ReLUSPMV
Benchmark
emiT
noitucexE
noitalumiS
8 8-opt 16 16-opt
Figure6.7: Executionperformancecomparisoncomparingthebaselineandimprovementsmadeto
theparallelengine.
Figure6.7showsacomparisonbetweenthesimulationexecutiontimeresultsofrunning
paralleltimingsimulationwith8and16cores. Thebaselineistheresultsrunningwiththeoriginal
implementation of the parallel simulation engine. We observe a slight execution time decrease,
indicatingthisstrategysuccessfullyreducestheblockingtime,butitdoesnothelpwiththescalability
issue. Toinvestigatefurther,wegeneratedanewsynchronizationblockingprofilegiventhisstrategy
applied. AportionoftheprofileisshowninFigure6.8. Thefullsynchronizationblockingprofile
canbefoundintheFigureA.3. Thisprofileverifiesthatthechannelblockingbottleneckhasbeen
eliminated. Thisexperimentdemonstratesthatthedesignforconcurrenteventexecution,alongwith
thesynchronizationmechanismplaysanimportantroleintheefficiencyoftheparallelsimulation
82

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
engine.
However, the profile suggests a different bottleneck. We can see that goroutines are
waiting to get the lock to execute the Send function of the Connection, which is a basic data
structuredefinedinAkita[81]. NaviSimleveragestheAkitasimulatorframeworkwhichprovidesa
port-connection-messagesystemtosupportcommunicationbetweencomponents,aswementioned
in section 4.2. The Akita simulator architecture offers a single connection to link multiple ports,
althougheachportcanonlybeservedbyasingleconnection. Whenconnectingmultipleports,the
connectionfunctionsasaround-robinarbitrator,managingboththesenderandreceiversidebuffers
ofeachport. Thiseliminatestheneedtodesignseparateswitchesandsimplifiestheconfiguration
process. However, it requires a mutual exclusion lock to protect the Send function, because the
Sendfunctionmanipulatesthebuffersofallports. WhenSendfunctionsarecalledfromdifferent
portsinmultiplegoroutinessimultaneously,theseoperationsmustbeserialized. Asignificantamount
ofblockingtimeoccurs,ascanbeseeninFigure6.8. Thisblockingtimeiscausedbytheconnection
thatconnectsallL2cachepartitionstotheL1scalar,instructionandvectorcachesfromallCompute
Units. TheR9NanoGPUconfiguration(seeTable6.2)has16L2cachepartitionsand16Shader
Arrays(eachhas4ComputeUnitswith6L1caches),resultingin96L1caches. Thefactthat16
writebackL2cachesaresendingmessagesto96L1cachesconcurrentlythroughasingleconnection
significantlyincreasesthechancethatmultiplegoroutinescalltheSendfunctionsimultaneously.
6.3 Performance Improvement Results
Basedonourprofilinganalysis,wehaveidentifiedseveralbottlenecksinbothemulation
and timing simulation. In this section, we apply software-level optimizations to mitigate these
bottlenecks to improve simulation speed. We focus on the emulation performance and show the
optimizationresults.
PerformanceDebuggingWeobserveReLUandPRtakemuchmoretimetofinishtheir
emulationthanotherbenchmarks. AftercheckingtheeventlistwhenexecutingReLU,wenoticea
largenumberofeventsscheduledandhandledintheComputeUnits. Wepinpointtheplacewhere
these events are created and find that most of the events are created due to the retry mechanism,
basicallyschedulingneweventstoretrysendingarequesttotheWorkgroupDispatcher. Thefailure
tosendtherequestcausestheschedulingofthenexteventtoretry.
Thisindicatesweencountertheretrymechanismpitfallwementionedinsubsection6.1.1.
We can avoid this by adjusting the way events are handled at the consumer end. Specifically for
83

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
Figure6.8: Thesynchronizationblockingprofileafterapplyingthereversing-orderstrategy.
84

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
NaviSim,wemodifyhowwehandleWGCompleteEventinourComputeUnitimplementation
during emulation. We use the original implementation of NaviSim as the baseline and we have
profiled the baseline code in subsection 6.2.2 to identify the performance issues. In the baseline
implementation,ifaworkgrouptotheWorkgroupDispatcheroperationfails,aretrywillbegenerated,
creating a new WGCompleteEvent. This retry will generate a new request to the Workgroup
Dispatcherandrepeatthisretryprocess. However,theWorkgroupDispatcher,whichisdesigned
to process one workgroup at a time struggles to process a large number of workgroups that are
completedwithinaveryshorttimeframe. Asaresult,thenumberofeventsandrequestsincreases
rapidly.
Toavoidnegativelyimpactingthetimingsimulation,whichcouldoccurifwemodifiedthe
WorkgroupDispatcher,weoptedtochangethewaywehandleWGCompleteEventintheCompute
Unit,giventhatbothemulationandtimingsimulationusethesamedispatcherimplementation. The
modifiedimplementationwaitsuntilallworkgroupsfinishandthensendsasinglerequest. Wealso
modifytheWorkgroupDispatcheraccordinglytohandlethisnewformofrequest.
Accordingtoourobservations(seesubsubsection6.2.2.2)wefoundthattheWorkgroup
Dispatcher does not run emulation efficiently, we propose and implement a new Workgroup Dis-
patcher tailored for emulation. The emulation dispatcher and timing simulation dispatcher share
thesameinterface,buthavetheirownimplementations. ThenewemulationWorkgroupDispatcher
removestiming-relatedmodelingandthisallowsforthesimultaneousdispatchofallworkgroupsin
akernel,alsoreducingthenumberofeventsusedfordispatching.
Figure6.9showstheexecutiontimeforthesimulatorforthethreecases: thebaseline,the
changethatmodifiedtheComputeUnit,andthenewWorkgroupDispatcherdedicatedtoemulation.
BothresolvetheemulationspeedissueforthePRandReLUbenchmarks,reducingtheexecutiontime
significantly. ThenewWorkgroupDispatcher,whichiscustom-designedforemulation,improvesthe
efficiencyofdispatchingworkgroupsforemulationandreducestheemulationtimeslightlyformost
benchmarks.
OverallPerformanceBasedontheprofilingresultswecollected,whichwereshownin
subsubsection6.2.2.1,weproposeanewmodeloftheALUthatdoesnotmodelthehardwarelogic
usedfortemporarilystoringoperandsfortheALU.Thisnewmodelmaintainsalowerlevelofdetail
inALU,allowingustoremovethescratchpaddatastructureandotherrelateddatastructures.
WereimplementsomecorefunctionsoftheALU,including: accessingoperandsfromtheregister
filesandeliminatingtheintermediatedatastructurescratchpadforstoringinputandoutputdata.
The redundant memory copy operations between buffers have been removed. We eliminate the
85

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
300
250
200
150
100
50
0
ATAX BS FWT FIR FLW KM MM MT PR ReLUSPMV
Benchmark
emiT
noitucexE
noitalumiS
Baseline Rectification Dispatcher
459.93 17726.33
Figure6.9: Performancecomparisonbetweenthebaseline,themodifiedComputeUnitandthenew
WorkgroupDispatcher.
bottleneckbyraisingthelevelofabstractioninmodelingtheALU.
We use the case where the new Workgroup Dispatcher is applied as our new baseline
simulator and provide a performance comparison between the original ALU and the new ALU
implementations. Theaveragespeedupgainis2.54.
The CPU profile generated by pprof indicates timing simulation suffers from a similar
performance issue caused by the use of the scratchpad and scratchpadpreparer data
structures. WemodifiedtheComputeUnitsbyutilizingthenewALUcomponentforemulationmode
optimization. The performance comparison and speedup are shown in Figure 6.11. The average
speedupis1.1. OurresultsshowthekernelexecutiontimepredictedbyNaviSimremainsaccurate
whenapplyingthenew,lessdetailedALUimplementation.
Insummary,weintroduceanewWorkgroupDispatcherforemulation,whichaddresses
akeysimulatorbottleneckandimprovesoverallperformance. WethenproposedanewALUmodelby
reimplementingtheinstructionemulationtoremovebothscratchpadandscratchpadpreparer
datastructures,whichinvolvedextensivecodemodificationbecauseeveryinstructionhastochange
thewaytheyaccessoperands. Applyingtheaboveoptimizations,wemanagedtoachieveanaverage
speedupof4.6×,reducingtheSlowdownfrom71,405to19,078,increasingtheKIPSby4×,and
loweringtheCIPGIby55.43%.
86

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
200
150
100
50
0
ATAX BS FWT FIR FLW KM MM MT PR ReLU SPMV
Benchmark
emiT
noitucexE
noitalumiS
Baseline New ALU
3
2
1
0
ATAX BS FWT FIR FLW KM MM MT PR ReLU SPMV
Benchmark
pudeepS
Figure 6.10: Emulation performance comparison between the original ALU and the new ALU
implementations.
87

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
800
600
400
200
0
ATAX BS FWT FIR FLW KM MM MT PR ReLU SPMV
Benchmark
emiT
noitucexE
noitalumiS
Baseline New ALU
1.2
1
0
ATAX BS FWT FIR FLW KM MM MT PR ReLU SPMV
Benchmark
pudeepS
Figure 6.11: Timing performance comparison between the original ALU and the new ALU
implementations.
88

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
ParallelScalabilityIssueTomitigatethesynchronizationblockingproblem,weexamined
thelogicusedforlockingthatprotectsshareddatastructuresinthecreateOrGetStorageUnit
functionandotherrelatedfunctions. Wearelookingforopportunitiestorelaxsomesynchroniza-
tion requirements, while maintaining correctness. The shared data structure is Storage, which
contains its addresses and corresponding StorageUnit. We find that every function call to
createOrGetStorageUnit is serialized, because most of the function code is treated as a
criticalsectionandisprotectedbyamutualexclusionlock. However,thisstrictformoflockingis
notnecessarywhenmultiplegoroutinesarereadingtheshareddatastructuresimultaneously.
SincefunctioncreateOrGetStorageUnitincludesbothcasesofreadingandwriting,
which achieves the initialization of the Storage, using a mutual exclusion lock is necessary to
ensurecorrectness. Torelaxthissynchronizationrequirement,wemodifiedtheimplementationto
haveafunctionnamedcreateStorageUnitforinitializationseparately,anditisguaranteedto
becalledbyasinglegoroutine. WehaveanotherfunctionnamedgetStorageUnittoreplace
thecreateOrGetStorageUnitandapplyareader/writermutualexclusionlock. Thelockcan
beheldbyanarbitrarynumberofreadersorasinglewriter. Wetakethesameapproachtoaddress
similarissuesinprocessTableandpageTableImpl(seeFigure6.4).
10
5
0
ATAX BS FWT FIR FLW KM MM MT PR ReLUSPMV
Benchmark
pudeepS
4 8 16
5.9
3.78
2.56
Geo
Mean
Figure6.12: Thespeedupachievedthroughparallelsimulationwhenrunningemulationmodeon
platform1.
Figure 6.12 shows the improved speedup achieved through parallel simulation when
running emulation mode. We normalize the simulation execution time of running with 4, 8, and
16coresagainstthetimeforasinglecore. Theaveragespeedupsfor8and16coresare3.78and
89

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
5.9, respectively. Compared to the original performance results (see Figure 6.2), we can see a
clearimprovementinthescalabilityissue. Figure6.12showstheimprovementwhenrunning16
coresgivesthebestperformance,whileintheoriginalcase,increasingto16coresdidnotimprove
scalability as compared to 8 cores.Although the speedup we achieve does not scale linearly with
thenumberofcores,weareabletoimprovetheparallelsimulationtimeandachievedanaverage
speedupof5.9×,aswellasamaximumspeedupof10.65×whenrunning16cores.
20
10
0
ATAX BS FWT FIR FLW KM MM MT PR ReLU SPMV
Benchmark
pudeepS
4 8 16 32 64
6.45
Geo
Mean
Figure6.13: Thespeedupachievedthroughparallelsimulationwhenrunningemulationmodeon
platform2.
Weevaluateouroptimizationsonadifferentmachinewith64cores(seeFigure6.13). The
averagespeedupsfor32and64coresare5.1and6.45,respectively. ExceptforATAX,ReLUand
SPMV,allbenchmarksexhibitanimprovementinspeedasweincreasethenumberofcores. With
64cores,weachievedanaveragespeedupof6.45×andamaximumspeedupof20.26×.
6.4 Summary
Throughoutourevaluationandperformanceimprovementresults,weelaborateonvarious
importantdesigndecisionsthatcouldleadtosignificantsimulationoverheadandimpacttheparallel
simulation scalability. Using NaviSim as an example, we apply our methodology for designing
high-performancesimulatorstocharacterizesomegeneralsimulatordesigndecisionsandidentifykey
factorsthatrestrictparallelsimulationscalability. Simulatordeveloperscanfollowourmethodology
whendesigningahigh-performancesimulatortominimizeunnecessarysimulationoverheadand
improve parallel simulation efficiency. We provide a detailed methodology, describing the steps
to design and optimize a simulator (See Table 6.4). This is an iterative process. After following
90

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
Table6.4: Stepstotunetheperformanceofasimulator
1 Evaluatethesimulationexecutionofthesimulatorinthemode(e.g.,emulation,timing,etc.)
youplantooptimizeandcollectanappropriatesetofperformancemetricsthatshouldinclude
executiontime,Slowdown,KIPSandCIPGIasthebaselinemeasurementset.
2 Profilethesimulationexecutionwithaperformanceprofilingtoolthatcanidentifyhotspots
wherethesimulationspendsthemosttime.
3 Identifythereasonwhythesimulationspendsthemosttimeinanyonefunctionorcomponent.
• Collecteventinformation(i.e.,function/componentname,thevirtualtimeinthesimu-
lator,theexecutiontime)foralleventsoreventsrelatedtothebottleneck. Inspectthem
intimeorder. Focusontheearliesttime-consumingevent.
• Revisitthemodelinglogic,lookingforopportunitiestoraisetheabstractionlevel,while
maintainsimulationaccuracy.
• Examineimplementationdetails(i.e.,datastructuresandbuilt-in/libraryfunctions)by
evaluatingthetimecomplexity.
4 Redesignormodifytherelatedpartsinthesimulatortomitigatethecauseofthebottleneck.
5 Rerun the simulation and collect the performance metrics. Compare the results with the
evaluationresultsfromStep1totesttheimpactofthechangesmadeinStep4.
ParallelSimulationTuning
1 Evaluatetheparallelsimulationscalabilityofthesimulatorwhenleveragingmultiplethreads.
Collecttheexecutiontimewhenmodifyingthenumberofthreads. Comparetheresultsand
analyzespeedupgainedthroughparallelism.
2 Profiletheparallelsimulationexecutionwithaperformanceprofilingtoolthatcanidentify
blockedthreadsduetosynchronizationoperations(i.e.,locksorsemaphores). Trytoutilized
advancedtoolsthatcanshowatraceofeachthread’sexecutionoversometime.
3 Review the event information at various timestamps in the simulator to ensure there are a
largeenoughnumberofeventstofullyexploitparallelexecution. Itisessentialtounderstand
howtheparallelsimulationenginemanageseventschedulingandhowconcurrenteventsare
handled. Makesurethreadsareusingsynchronizationproperlytoensurethecorrectorderof
operation.
4 Inspectaccessestosharedresourcesperformedbyeachthread,identifyingresourcesthatare
accessedheavilyduringsimulation.
5 Relaxthesynchronizationrequirementsifpossible. Ifwecannotrelaxtheserequirements,
thenconsiderrearchitectingthemanagementofsharedresourcesinadistributedmanner.
6 Reevaluatetheparallelsimulationscalabilityofthesimulatortotesttheimpactofthechanges
madeinStep5.
91

CHAPTER6. SIMULATIONPERFORMANCEEVALUATIONANDOPTIMIZATION
thesesteps,developerscanrevisitprofilingresultsandcontinuerefiningthesimulatorbasedonthe
observedbehavior. WedemonstrateathoroughperformanceanalysisofNaviSimbyprofilingboth
emulation and timing simulation execution running in serial and parallel. Based on the profiling
results,wefocusedonimprovingemulationperformancebyproposingnewmodelsfortheWorkgroup
DispatcherandALU.Wealsomitigatedthesynchronizationblockingbottlenecks,improvingthe
parallelemulationscalability.
92

Chapter 7
Conclusion and Future Work
7.1 Conclusion
This dissertation has focused on addressing the challenges involved when designing
an architectural simulation framework that features high adaptability and high performance. We
have developed a novel design methodology for producing a highly adaptable simulator. We
have introduced a methodology for examining the key simulator design decisions that impact
the simulation performance. We presented NaviSim, a AMD GPU simulator that was designed,
implemented,andvalidated,followingourprinciples. NaviSimisabletoaccommodatenewfeatures
and design extensions, allowing the user to explore a wide range of GPU architecture questions.
WeperformperformanceanalysisandtuningonNaviSim,followingourmethodologyofdesigning
high-performancesimulators. Wesummarizethecontributionsofthisdissertationasfollows:
• Wepresentthedesigngoalsanddevelopmentmethodologyforbuildingahighlyadaptable
and accurate GPU simulator. To build a GPU simulator that can support multiple ISAs
andmicroarchitectures,andbeeasilyextendedtosupportfutureincomingarchitectures,we
demonstrate our approach by presenting NaviSim [21]. NaviSim models the state-of-the-
artAMDRDNAandCDNAarchitectures,providingbothhighmodularityandextensibility.
Throughourintensiveandrigorousvalidationmethodology,whichincludeddevelopingseveral
microbenchmarks,NaviSimachievesasmall9.85%averagesimulationerror,ascomparedto
hardwareexecution. Wefollowedoursimulatordevelopmentmethodology. Ourapproachcan
begeneralizedtodesignanyarchitecturalsimulator,producingatoolsetdesignedwithhigh
modularity,extensibility,andhighfidelity. Weconductedaperformancestudyoftheimpact
93

CHAPTER7. CONCLUSIONANDFUTUREWORK
ofindividualarchitecturefeaturesofAMDRNDAarchitecture,showcasingthehighflexibility
andconfigurabilityofNaviSim.
• WepresentaseriesofGPU-basedmicroarchitecturalextensions(GME)[69]asausecaseof
designspaceexplorationusingNaviSim. Wedemonstratethehighdegreeofflexibilityand
adaptabilityofNaviSim. First,NaviSimservesasaperformanceanalysistoolbyidentifying
themicroarchitecture-levelbottlenecks. Next,wedemonstratethehighflexibilityofNaviSim,
whichenablesuserstoquicklyexplorethedesignspaceofnewarchitecturalfeaturesoffuture
GPUs. UserscaneasilyextendNaviSimwithnewcomponents(e.g.,anewon-chipnetwork
and a new modular operations unit), and connect them to the existing components without
modifyingtherestofthesimulator. WewereabletomarryNaviSimtoBlockSim,allowingus
toperformhybridsimulation. ThisusecasedemonstratesthehighconfigurablityofNaviSim.
Architecturalparameters,suchasthesizeofLDSandregisterfiles,canbechangedeffortlessly.
Last but not least, NaviSim, as an execution-driven simulator, can support the decode and
executionofnewlydesignedinstructions,whichenablesthefullsupportofGPUarchitecture
design,includingbothISAandmicroarchitecture.
• Wepresentamethodologyforexaminingthekeysimulatordesigndecisionsthatsimulator
developersneedtoconsidertominimizeunnecessarysimulationoverheadandimproveparallel
simulationefficiency. Weintroduceasetofperformancemetricsforthearchitecturalsimulation
community,sothattheycanutilizestandardmetricswhenevaluatingtheexecutionperformance
ofanyarchitecturalsimulator. UsingNaviSimasanexample,weappliedourmethodologyto
thedesignofourhigh-performancesimulator. Wewereabletoprovideguidanceonsimulator
design decisions and identify key factors that impact parallel simulation scalability. Our
methodologycanbegeneralizedtoshapethedesignprocessofanyarchitecturalsimulatorin
ordertoachievehigh-performance. Next,weconductedathoroughperformanceanalysisof
NaviSim to find bottlenecks and performance issues. We built a process that leverages the
executionprofileofNaviSimtohelpguideasetofsoftware-leveloptimizationstoimprove
simulationspeed.
7.2 Future Work
The main goal of this dissertation was to propose a design methodology for producing
highly-adaptableandhigh-performancesimulationframeworks. WeleveragedtheAkitasimulator
94

CHAPTER7. CONCLUSIONANDFUTUREWORK
frameworktoenablemodularandextensibledesignsaswellasparallelsimulationtoachievehigh
performance. Akita has a rich set of features, one of which is the support of parallel simulation.
Akitaemploysaparallelsimulationenginetoachieveparallelsimulation. Throughourevaluationof
runningNaviSimonmultiplecoresandfurtherinvestigationofsimulationscalability,wefoundthat
Akitahasknownlimitationsinitsparallelsimulationsupport. Onefruitfuldirectionforfuturework
wouldbetoimprovetheparallelsimulationsupportofAkita.
AkeylimitationinAkitascalabilityisrootedinitsport-connection-messagesystem,which
isoneofAkita’scoreconceptstoenablethecommunicationbetweencomponents,whilekeeping
eachcomponentfullyencapsulated. Thistypeofuniversalcommunicationsystemofconnecting
allcomponentsleadstosignificantsynchronizationoverhead,asdifferentmessagesinterferewith
eachotherbysharingthesameconnection. Wewouldneedtoredesignthecommunicationsystem
to reduce the synchronization overhead. To achieve better parallel simulation scalability, further
opportunitiescanbeexploitedbyredesigningtheparallelsimulationenginetohaveamoreaggressive
wayofmanagingeventschedulingandutilizationofmultiplethreads.
95

Bibliography
[1] L.-N.Pouchetetal.,“Polybench: Thepolyhedralbenchmarksuite,”2012.
[2] A. Inc., “Amd graphics core next architecture,” url https://www.techpowerup.com/gpu-
specs/docs/amd-gcn1-architecture.pdf,2012.
[3] AMD Inc., “Introducing rdna architecture, the all new radeon gaming architecture
powering“navi”,”2019.[Online].Available: https://www.amd.com/system/files/documents/
rdna-whitepaper.pdf
[4] Y.Lecun,L.Bottou,Y.Bengio,andP.Haffner,“Gradient-basedlearningappliedtodocument
recognition,”ProceedingsoftheIEEE,vol.86,no.11,pp.2278–2324,1998.
[5] W.Liu,Z.Wang,X.Liu,N.Zeng,Y.Liu,andF.E.Alsaadi,“Asurveyofdeepneuralnetwork
architectures and their applications,” Neurocomputing, vol. 234, pp. 11–26, 2017. [Online].
Available: https://www.sciencedirect.com/science/article/pii/S0925231216315533
[6] Top500,https://www.top500.org/lists/top500/.
[7] NVIDIA,“Computeunifieddevicearchitectureprogrammingguide,”2007.
[8] A.Inc.,“Hipprogrammingguide,”https://rocmdocs.amd.com/en/latest/Programming Guides/
HIP-GUIDE.html,2022.
[9] Y.Sun,T.Baruah,andD.Kaeli,“Acceleratedcomputingwithhip,”2023.
[10] D.R.Kaeli,P.Mistry,D.Schaa,andD.P.Zhang,HeterogeneouscomputingwithOpenCL2.0.
Burlington,MA,USA:MorganKaufmann,2015.
[11] T.Nowatzki,J.Menon,C.-H.Ho,andK.Sankaralingam,“Architecturalsimulatorsconsidered
harmful,”IEEEMicro,vol.35,no.6,pp.4–12,2015.
96

BIBLIOGRAPHY
[12] A.Inc.,“Introducingcdnaarchitecture,theall-newAMDGPUarchitectureforthemodernera
ofhpc&ai,”https://www.amd.com/system/files/documents/amd-cdna-whitepaper.pdf,2020.
[13] J. H. Cheon, A. Kim, M. Kim, and Y. Song, “Homomorphic encryption for arithmetic of
approximate numbers,” in Advances in Cryptology–ASIACRYPT 2017: 23rd International
Conference on the Theory and Applications of Cryptology and Information Security, Hong
Kong,China,December3-7,2017,Proceedings,PartI23. Springer,2017,pp.409–437.
[14] I. Chillotti, N. Gama, M. Georgieva, and M. Izabache`ne, “Tfhe: Fast fully homomorphic
encryptionoverthetorus,”JournalofCryptology,vol.33,042019.
[15] W.Jung,E.Lee,S.Kim,J.Kim,N.Kim,K.Lee,C.Min,J.H.Cheon,andJ.H.Ahn,“Acceler-
atingfullyhomomorphicencryptionthrougharchitecture-centricanalysisandoptimization,”
IEEEAccess,vol.9,pp.98772–98789,2021.
[16] K.Shivdikar,G.Jonatan,E.Mora,N.Livesay,R.Agrawal,A.Joshi,J.L.Abella´n,J.Kim,and
D.Kaeli,“AcceleratingpolynomialmultiplicationforhomomorphicencryptiononGPUs,”in
2022IEEEInternationalSymposiumonSecureandPrivateExecutionEnvironmentDesign
(SEED). IEEE,2022,pp.61–72.
[17] N.Livesay,G.Jonatan,E.Mora,K.Shivdikar,R.Agrawal,A.Joshi,J.L.Abella´n,J.Kim,and
D.Kaeli,“AcceleratingfinitefieldarithmeticforhomomorphicencryptiononGPUs,”2023
IEEEMICRO,2023.
[18] A.AlBadawi,B.Veeravalli,J.Lin,N.Xiao,M.Kazuaki,andA.K.M.Mi,“Multi-GPUdesign
andperformanceevaluationofhomomorphicencryptiononGPUclusters,”IEEETransactions
onParallelandDistributedSystems,vol.32,no.2,pp.379–391,2020.
[19] S. Williams, A. Waterman, and D. Patterson, “Roofline: An insightful visual performance
modelformulticorearchitectures,”Commun.ACM,vol.52,no.4,p.65–76,apr2009.[Online].
Available: https://doi.org/10.1145/1498765.1498785
[20] Y. Sun, T. Baruah, S. A. Mojumder, S. Dong, X. Gong, S. Treadway, Y. Bao, S. Hance,
C. McCardwell, V. Zhao, H. Barclay, A. K. Ziabari, Z. Chen, R. Ubal, J. L. Abella´n,
J. Kim, A. Joshi, and D. Kaeli, “Mgpusim: Enabling multi-gpu performance modeling and
optimization,”inProceedingsofthe46thInternationalSymposiumonComputerArchitecture,
97

BIBLIOGRAPHY
ser. ISCA ’19. New York, NY, USA: Association for Computing Machinery, 2019, p.
197–209.[Online].Available: https://doi.org/10.1145/3307650.3322230
[21] Y.Bao,Y.Sun,Z.Feric,M.T.Shen,M.Weston,J.L.Abella´n,T.Baruah,J.Kim,A.Joshi,
andD.Kaeli,“Navisim: Ahighlyaccurategpusimulatorforamdrdnagpus,”inProceedings
of the International Conference on Parallel Architectures and Compilation Techniques, ser.
PACT’22. NewYork,NY,USA:AssociationforComputingMachinery,2023,p.333–345.
[Online].Available: https://doi.org/10.1145/3559009.3569666
[22] J.J.Y.TimothySherwood,“Computerarchitecturesimulationandmodeling,”IEEEComputer
Society,vol.26,pp.5–7,July/August2006.
[23] D. Sanchez and C. Kozyrakis, “Zsim: Fast and accurate microarchitectural simulation of
thousand-core systems,” in Proceedings of the 40th Annual International Symposium on
Computer Architecture, ser. ISCA ’13. New York, NY, USA: Association for Computing
Machinery,2013,p.475–486.[Online].Available: https://doi.org/10.1145/2485922.2485963
[24] A. Mohammad, U. Darbaz, G. Dozsa, S. Diestelhorst, D. Kim, and N. S. Kim, “dist-gem5:
Distributed simulation of computer clusters,” in 2017 IEEE International Symposium on
PerformanceAnalysisofSystemsandSoftware(ISPASS),2017,pp.153–162.
[25] R. Wunderlich, T. Wenisch, B. Falsafi, and J. Hoe, “Smarts: accelerating microarchitecture
simulation via rigorous statistical sampling,” in 30th Annual International Symposium on
ComputerArchitecture,2003.Proceedings.,2003,pp.84–95.
[26] C.AvalosBaddouh, M.Khairy, R.N.Green, M.Payer, andT.G.Rogers, “Principalkernel
analysis: A tractable methodology to simulate scaled gpu workloads,” in MICRO-54: 54th
AnnualIEEE/ACMInternationalSymposiumonMicroarchitecture,ser.MICRO’21. New
York,NY,USA:AssociationforComputingMachinery,2021,p.724–737.[Online].Available:
https://doi.org/10.1145/3466752.3480100
[27] C.Liu,Y.Sun,andT.E.Carlson,“Photon: Afine-grainedsampledsimulationmethodology
forgpuworkloads,”inProceedingsofthe56thAnnualIEEE/ACMInternationalSymposium
on Microarchitecture, ser. MICRO ’23. New York, NY, USA: Association for Computing
Machinery,2023,p.1227–1241.[Online].Available: https://doi.org/10.1145/3613424.3623773
98

BIBLIOGRAPHY
[28] O.Villa,D.Lustig,Z.Yan,E.Bolotin,Y.Fu,N.Chatterjee,N.Jiang,andD.Nellans,“Need
for speed: Experiences building a trustworthy system-level gpu simulator,” in 2021 IEEE
InternationalSymposiumonHigh-PerformanceComputerArchitecture(HPCA),IEEE. Seoul,
Korea(South): IEEE,2021,pp.868–880.
[29] M. Auguin, F. Boeri, J. Dalban, and A. Vincent-Carrefour, “Experience using a simd/spmd
multiprocessor architecture,” Microprocessing and Microprogramming, vol. 21, no. 1, pp.
171–177, 1987, microcomputers: Usage, Methods and Structures. [Online]. Available:
https://www.sciencedirect.com/science/article/pii/0165607487900342
[30] J.KimandY.Kim,“Hbm: Memorysolutionforbandwidth-hungryprocessors,”in2014IEEE
HotChips26Symposium(HCS),IEEE. Cupertino,CA:IEEE,2014,pp.1–24.
[31] J.JESD250,“Graphicsdoubledatarate6(gddr6)sgramstandard,”JEDECSolidStateTech-
nologyAssociation,2017.
[32] A.AkramandL.Sawalha,“Asurveyofcomputerarchitecturesimulationtechniquesandtools,”
IEEEAccess,vol.7,pp.78120–78145,2019.
[33] J. Banks, J. S. Carson, B. L. Nelson, and D. M. Nicol, Discrete-Event System Simulation.
PrenticeHallPublishing,1995.[Online].Available: https://api.semanticscholar.org/CorpusID:
122566976
[34] R.M.Fujimoto,“Paralleldiscreteeventsimulation,”Commun.ACM,vol.33,no.10,p.30–53,
oct1990.[Online].Available: https://doi.org/10.1145/84537.84545
[35] C.Dufour,C.Andrade,andJ.Be´langer,“Real-timesimulationtechnologiesineducation: a
linktomodernengineeringmethodsandpractices,”inProceedingsofthe11thInternational
ConferenceonEngineeringandTechnologyEducationINTERTECH,2010.
[36] T.GleixnerandI.Molnar,“perf,”2008.[Online].Available: https://github.com/torvalds/linux
[37] S.L.Graham,P.B.Kessler,andM.K.Mckusick,“Gprof: Acallgraphexecutionprofiler,”in
Proceedingsofthe1982SIGPLANSymposiumonCompilerConstruction,ser.SIGPLAN’82.
New York, NY, USA: Association for Computing Machinery, 1982, p. 120–126. [Online].
Available: https://doi.org/10.1145/800230.806987
[38] Microsoft, “Windows performance toolkit,” 2022. [Online]. Available: https://github.com/
microsoft/microsoft-performance-toolkit-sdk
99

BIBLIOGRAPHY
[39] Google,“pprof.”[Online].Available: https://github.com/google/pprof
[40] I.C.Laboratory,“Papi,”2009.[Online].Available: https://icl.utk.edu/papi/
[41] P. Paillier, “Public-key cryptosystems based on composite degree residuosity classes,” in
Advances in Cryptology — EUROCRYPT ’99, J. Stern, Ed. Berlin, Heidelberg: Springer
BerlinHeidelberg,1999,pp.223–238.
[42] R. L. Rivest, A. Shamir, and L. Adleman, “A method for obtaining digital signatures and
public-key cryptosystems,” Commun. ACM, vol. 21, no. 2, p. 120–126, feb 1978. [Online].
Available: https://doi.org/10.1145/359340.359342
[43] C.Gentry,“Fullyhomomorphicencryptionusingideallattices,”inProceedingsoftheforty-first
annualACMsymposiumonTheoryofcomputing,2009,pp.169–178.
[44] D.MicciancioandO.Regev,“Lattice-basedcryptography,”Post-quantumcryptography,pp.
147–191,2009.
[45] V.Shoup,Acomputationalintroductiontonumbertheoryandalgebra. CambridgeUniversity
Press,2009.
[46] A.J.RashmiAgrawal,OnArchitectingFullyHomomorphicEncryption-basedComputingSys-
tems. SpringerCham,2023.[Online].Available: https://doi.org/10.1007/978-3-031-31754-5
[47] Y.Fu,E.Bolotin,N.Chatterjee,D.Nellans,andS.W.Keckler,“Gpudomainspecializationvia
composable on-package architecture,” ACM Trans. Archit. Code Optim., vol. 19, no. 1, dec
2021.[Online].Available: https://doi.org/10.1145/3484505
[48] N. Binkert, B. Beckmann, G. Black, S. K. Reinhardt, A. Saidi, A. Basu, J. Hestness, D. R.
Hower,T.Krishna,S.Sardashti,R.Sen,K.Sewell,M.Shoaib,N.Vaish,M.D.Hill,andD.A.
Wood, “The gem5 simulator,” SIGARCH Comput. Archit. News, vol. 39, no. 2, p. 1–7, aug
2011.[Online].Available: https://doi.org/10.1145/2024716.2024718
[49] A.Gutierrez,B.M.Beckmann,A.Dutu,J.Gross,M.LeBeane,J.Kalamatianos,O.Kayiran,
M.Poremba,B.Potter,S.Puthooretal.,“Lostinabstraction: Pitfallsofanalyzinggpusatthe
intermediatelanguagelevel,”in2018IEEEInternationalSymposiumonHighPerformance
ComputerArchitecture(HPCA),IEEE. Vienna,Austria: IEEE,2018,pp.608–619.
100

BIBLIOGRAPHY
[50] C.Collange,M.Daumas,D.Defour,andD.Parello,“Barra: Aparallelfunctionalsimulator
forgpgpu,”in2010IEEEInternationalSymposiumonModeling,AnalysisandSimulationof
ComputerandTelecommunicationSystems. Miami,Florida,USA:IEEE,2010,pp.351–360.
[51] N.Farooqui,A.Kerr,G.Diamos,S.Yalamanchili,andK.Schwan,“Aframeworkfordynami-
callyinstrumentinggpucomputeapplicationswithingpuocelot,”inProceedingsoftheFourth
WorkshoponGeneralPurposeProcessingonGraphicsProcessingUnits. NewportBeach,
CA:ACM,2011,pp.1–9.
[52] A. Kerr, G. Diamos, and S. Yalamanchili, “Gpu application development, debugging, and
performancetuningwithgpuocelot,”inGPUComputingGemsJadeEdition. Amsterdam,
Netherlands: Elsevier,2012,pp.409–427.
[53] A.Bakhoda,G.L.Yuan,W.W.Fung,H.Wong,andT.M.Aamodt,“Analyzingcudaworkloads
using a detailed gpu simulator,” in 2009 IEEE International Symposium on Performance
AnalysisofSystemsandSoftware,IEEE. Boston,MAUSA:IEEE,2009,pp.163–174.
[54] H.Wang,F.Luo,M.Ibrahim,O.Kayiran,andA.Jog,“Efficientandfairmulti-programming
in gpus via effective bandwidth management,” in 2018 IEEE International Symposium on
HighPerformanceComputerArchitecture(HPCA),IEEE. Vienna,Austria: IEEE,2018,pp.
247–258.
[55] R.Ausavarungnirun,J.Landgraf,V.Miller,S.Ghose,J.Gandhi,C.J.Rossbach,andO.Mutlu,
“Mosaic: agpumemorymanagerwithapplication-transparentsupportformultiplepagesizes,”
inProceedingsofthe50thAnnualIEEE/ACMInternationalSymposiumonMicroarchitecture.
Cambridge,MA,USA:IEEE,2017,pp.136–150.
[56] M. Khairy, Z. Shen, T. M. Aamodt, and T. G. Rogers, “Accel-sim: An extensible simula-
tionframeworkforvalidatedgpumodeling,”in2020ACM/IEEE47thAnnualInternational
Symposium on Computer Architecture (ISCA), IEEE. Valencia, Spain: IEEE, 2020, pp.
473–486.
[57] O. Villa, M. Stephenson, D. Nellans, and S. W. Keckler, “Nvbit: A dynamic binary
instrumentationframeworkfornvidiagpus,”inProceedingsofthe52ndAnnualIEEE/ACM
International Symposium on Microarchitecture, ser. MICRO ’52. New York, NY,
USA: Association for Computing Machinery, 2019, p. 372–383. [Online]. Available:
https://doi.org/10.1145/3352460.3358307
101

BIBLIOGRAPHY
[58] R. Ubal, B. Jang, P. Mistry, D. Schaa, and D. Kaeli, “Multi2sim: A simulation
framework for cpu-gpu computing,” in Proceedings of the 21st International Conference
on Parallel Architectures and Compilation Techniques, ser. PACT ’12. New York, NY,
USA: Association for Computing Machinery, 2012, p. 335–344. [Online]. Available:
https://doi.org/10.1145/2370816.2370865
[59] X.Gong,R.Ubal,andD.Kaeli,“Multi2simkepler: Adetailedarchitecturalgpusimulator,”
in 2017 IEEE International Symposium on Performance Analysis of Systems and Software
(ISPASS),IEEE. SantaRosa,CA:IEEE,2017,pp.269–278.
[60] H. Kim, J. Lee, N. B. Lakshminarayana, J. Sim, J. Lim, and T. Pho, Macsim: A cpu-gpu
heterogeneoussimulationframeworkuserguide,GeorgiaInstituteofTechnology,Atlanta,GA,
2012.
[61] G.Malhotra,S.Goel,andS.R.Sarangi,“Gputejas: Aparallelsimulatorforgpuarchitectures,”
in 2014 21st International Conference on High Performance Computing (HiPC), 2014, pp.
1–10.
[62] S.LeeandW.W.Ro,“Parallelgpuarchitecturesimulationframeworkexploitingworkallo-
cationunitparallelism,”in2013IEEEInternationalSymposiumonPerformanceAnalysisof
SystemsandSoftware(ISPASS),2013,pp.107–117.
[63] ——,“Parallelgpuarchitecturesimulationframeworkexploitingarchitectural-levelparallelism
withtimingerrorprediction,”IEEETransactionsonComputers,vol.65,no.4,pp.1253–1265,
2016.
[64] R. Murphy, A. F. Rodrigues, P. Kogge, and K. D. Underwood, “The structural simulation
toolkit : A tool for bridging the architectural/microarchitectural evaluation gap.” 12 2004.
[Online].Available: https://www.osti.gov/biblio/1088092
[65] C. Hughes, S. D. Hammond, M. Zhang, Y. Liu, T. Rogers, and R. J. Hoekstra, “Sst-gpu:
A scalable sst gpu component for performance modeling and profiling.” Sandia National
Lab.(SNL-NM),Albuquerque,NM(UnitedStates),Tech.Rep.,2021.
[66] E.C.Marangoz,K.-D.Kang,andS.Shin,“Designinggpuarchitectureformemorybandwidth
reservation,”in2021IEEEInternationalSymposiumonPerformanceAnalysisofSystemsand
Software(ISPASS),2021,pp.87–89.
102

BIBLIOGRAPHY
[67] A.Arunkumar,E.Bolotin,B.Cho,U.Milic,E.Ebrahimi,O.Villa,A.Jaleel,C.-J.Wu,and
D. Nellans, “Mcm-gpu: Multi-chip-module gpus for continued performance scalability,” in
2017 ACM/IEEE 44th Annual International Symposium on Computer Architecture (ISCA),
2017,pp.320–332.
[68] G.-J. V. D. Braak and H. Corporaal, “R-gpu: A reconfigurable gpu architecture,”
ACM Trans. Archit. Code Optim., vol. 13, no. 1, mar 2016. [Online]. Available:
https://doi.org/10.1145/2890506
[69] K. Shivdikar, Y. Bao, R. Agrawal, M. Shen, G. Jonatan, E. Mora, A. Ingare, N. Livesay,
J. L. AbellA´N, J. Kim, A. Joshi, and D. Kaeli, “Gme: Gpu-based microarchitectural
extensions to accelerate homomorphic encryption,” in Proceedings of the 56th Annual
IEEE/ACM International Symposium on Microarchitecture, ser. MICRO ’23. New York,
NY, USA: Association for Computing Machinery, 2023, p. 670–684. [Online]. Available:
https://doi.org/10.1145/3613424.3614279
[70] A.AlBadawi,L.Hoang,C.F.Mun,K.Laine,andK.M.M.Aung,“Privft: Privateandfasttext
classificationwithhomomorphicencryption,”IEEEAccess,vol.8,pp.226544–226556,2020.
[71] A.AlBadawi,B.Veeravalli,C.F.Mun,andK.M.M.Aung,“High-performancefvsomewhat
homomorphicencryptiononGPUs: Animplementationusingcuda,”IACRTransactionson
CryptographicHardwareandEmbeddedSystems,pp.70–95,2018.
[72] W.Jung,S.Kim,J.H.Ahn,J.H.Cheon,andY.Lee,“Over100xfasterbootstrappinginfully
homomorphicencryptionthroughmemory-centricoptimizationwithGPUs,”IACRTransactions
onCryptographicHardwareandEmbeddedSystems,pp.114–148,2021.
[73] M. S. Riazi, K. Laine, B. Pelton, and W. Dai, “Heax: An architecture for computing on
encrypteddata,”inProceedingsoftheTwenty-FifthInternationalConferenceonArchitectural
SupportforProgrammingLanguagesandOperatingSystems,2020,pp.1295–1309.
[74] S.Kim, K.Lee, W.Cho, Y.Nam, J.H.Cheon, andR.A.Rutenbar, “Hardwarearchitecture
of a number theoretic transform for a bootstrappable RNS-based homomorphic encryption
scheme,”in2020IEEE28thAnnualInternationalSymposiumonField-ProgrammableCustom
ComputingMachines(FCCM). IEEE,2020,pp.56–64.
103

BIBLIOGRAPHY
[75] R. Agrawal, L. de Castro, G. Yang, C. Juvekar, R. Yazicigil, A. Chandrakasan, V. Vaikun-
tanathan,andA.Joshi,“Fab: Anfpga-basedacceleratorforbootstrappablefullyhomomorphic
encryption,”in2023IEEEInternationalSymposiumonHigh-PerformanceComputerArchitec-
ture(HPCA). IEEE,2023,pp.882–895.
[76] N.Samardzic,A.Feldmann,A.Krastev,S.Devadas,R.Dreslinski,C.Peikert,andD.Sanchez,
“F1: Afastandprogrammableacceleratorforfullyhomomorphicencryption,”inMICRO-54:
54thAnnualIEEE/ACMInternationalSymposiumonMicroarchitecture,2021,pp.238–252.
[77] S.Kim,J.Kim,M.J.Kim,W.Jung,J.Kim,M.Rhu,andJ.H.Ahn,“BTS:Anacceleratorfor
bootstrappablefullyhomomorphicencryption,”inProceedingsofthe49thAnnualInternational
SymposiumonComputerArchitecture,2022,pp.711–725.
[78] N.Samardzic,A.Feldmann,A.Krastev,N.Manohar,N.Genise,S.Devadas,K.Eldefrawy,
C. Peikert, and D. Sanchez, “Craterlake: a hardware accelerator for efficient unbounded
computationonencrypteddata,”inProceedingsofthe49thAnnualInternationalSymposium
onComputerArchitecture,2022,pp.173–187.
[79] J.Kim,G.Lee,S.Kim,G.Sohn,M.Rhu,J.Kim,andJ.H.Ahn,“Ark: Fullyhomomorphic
encryptionacceleratorwithruntimedatagenerationandinter-operationkeyreuse,” in2022
55thIEEE/ACMInternationalSymposiumonMicroarchitecture(MICRO). IEEE,2022,pp.
1237–1254.
[80] R.C.Martin,AgileSoftwareDevelopment: Principles,Patterns,andPractices. USA:Prentice
HallPTR,2003.
[81] Y.Sun,EnablingCollaborativeHeterogeneousComputing. NortheasternUniversity,2020.
[Online].Available: https://books.google.com/books?id=dMeqzgEACAAJ
[82] V. Volkov and J. W. Demmel, “Benchmarking gpus to tune dense linear algebra,” in SC’08:
Proceedingsofthe2008ACM/IEEEconferenceonSupercomputing,IEEE. Austin,TX,USA:
IEEE,2008,pp.1–11.
[83] T.-B. Chiou, A. C. Chen, M. Dusa, and S.-E. Tseng, “Impact of euv patterning scenario on
differentdesignstylesandtheirgroundrulesfor7nm/5nmnodebeollayers,”inDesign-Process-
TechnologyCo-optimizationforManufacturabilityX,vol.9781,InternationalSocietyforOptics
andPhotonics. Bellingham,WashingtonUSA:SPIE,2016,p.978107.
104

BIBLIOGRAPHY
[84] Y. Sun, S. Mukherjee, T. Baruah, S. Dong, J. Gutierrez, P. Mohan, and D. Kaeli,
“Evaluating performance tradeoffs on the radeon open compute platform,” in 2018 IEEE
InternationalSymposiumonPerformanceAnalysisofSystemsandSoftware(ISPASS). Los
Alamitos, CA, USA: IEEE Computer Society, apr 2018, pp. 209–218. [Online]. Available:
https://doi.ieeecomputersociety.org/10.1109/ISPASS.2018.00034
[85] J.J.K.Park,Y.Park,andS.Mahlke,“Dynamicresourcemanagementforefficientutilization
of multitasking gpus,” in Proceedings of the Twenty-Second International Conference on
ArchitecturalSupportforProgrammingLanguagesandOperatingSystems,ser.ASPLOS’17.
New York, NY, USA: Association for Computing Machinery, 2017, p. 527–540. [Online].
Available: https://doi.org/10.1145/3037697.3037707
[86] L.Wang,M.Huang,andT.El-Ghazawi,“Exploitingconcurrentkernelexecutionongraphic
processing units,” in 2011 International Conference on High Performance Computing &
Simulation,IEEE. Istanbul,Turkey: IEEE,2011,pp.24–32.
[87] S.Li,Z.Yang,D.Reddy,A.Srivastava,andB.Jacob,“Dramsim3: acycle-accurate,thermal-
capable dram simulator,” IEEE Computer Architecture Letters, vol. 19, no. 2, pp. 106–109,
2020.
[88] A.Staff,“Openclandtheamdappsdkv2.4,”2014.
[89] A.Danalis,G.Marin,C.McCurdy,J.S.Meredith,P.C.Roth,K.Spafford,V.Tipparaju,and
J.S.Vetter,“Thescalableheterogeneouscomputing(shoc)benchmarksuite,”inProceedings
of the 3rd Workshop on General-Purpose Computation on Graphics Processing Units, ser.
GPGPU-3. New York, NY, USA: Association for Computing Machinery, 2010, p. 63–74.
[Online].Available: https://doi.org/10.1145/1735688.1735702
[90] Y. Sun, X. Gong, A. Ziabari, L. Yu, X. Li, S. Mukherjee, C. Mccardwell, A. Villegas,
and D. Kaeli, “Hetero-mark, a benchmark suite for cpu-gpu collaborative computing,”
in 2016 IEEE International Symposium on Workload Characterization (IISWC). Los
Alamitos, CA, USA: IEEE Computer Society, sep 2016, pp. 1–10. [Online]. Available:
https://doi.ieeecomputersociety.org/10.1109/IISWC.2016.7581262
[91] S. Dong and D. Kaeli, “Dnnmark: A deep neural network benchmark suite for
gpus,” in Proceedings of the General Purpose GPUs, ser. GPGPU-10. New York,
105

BIBLIOGRAPHY
NY, USA: Association for Computing Machinery, 2017, p. 63–72. [Online]. Available:
https://doi.org/10.1145/3038228.3038239
[92] H.A.VanderVorst,“Bi-cgstab: Afastandsmoothlyconvergingvariantofbi-cgforthesolution
ofnonsymmetriclinearsystems,”SIAMJournalonscientificandStatisticalComputing,vol.13,
no.2,pp.631–644,1992.
[93] Y.Sun,Y.Zhang,A.Mosallaei,M.D.Shah,C.Dunne,andD.Kaeli,“Daisen: Aframework
forvisualizingdetailedGPUexecution,”EurographicsConferenceonVisualization,vol.40,
no.3,pp.239–250,2021.
[94] N.E.Jerger,T.Krishna,andL.-S.Peh,“On-chipnetworks,”SynthesisLecturesonComputer
Architecture,vol.12,no.3,pp.1–210,2017.
[95] J.BalfourandW.J.Dally,“Designtradeoffsfortiledcmpon-chipnetworks,”inACMInterna-
tionalconferenceonsupercomputing25thanniversaryvolume,2006,pp.390–401.
[96] C.WalshawandM.Cross,“Multilevelmeshpartitioningforheterogeneouscommunication
networks,”Futuregenerationcomputersystems,vol.17,no.5,pp.601–623,2001.
[97] K.Han,S.Hong,J.H.Cheon,andD.Park,“Logisticregressiononhomomorphicencrypted
dataatscale,”inProceedingsoftheAAAIconferenceonartificialintelligence,vol.33,no.01,
2019,pp.9466–9471.
[98] E.Lee,J.-W.Lee,J.Lee,Y.-S.Kim,Y.Kim,J.-S.No,andW.Choi,“Low-complexitydeep
convolutionalneuralnetworksonfullyhomomorphicencryptionusingmultiplexedparallel
convolutions,”inInternationalConferenceonMachineLearning. PMLR,2022,pp.12403–
12422.
[99] R. M. Fujimoto, “Parallel discrete event simulation,” Communications of the ACM, vol. 33,
no.10,pp.30–53,1990.
106

Appendix A
| Full Synchronization | Blocking | Profiles |
| -------------------- | -------- | -------- |
107

APPENDIXA. FULLSYNCHRONIZATIONBLOCKINGPROFILES
File: relu
Build ID: 3f5eed249fe4f7c03971670e0927d11d67b652a6
Type: delay
| Time: Jul 2, 2024 at 5:42pm (EDT)                           |     |     | sim                             | driver              | driver             |
| ----------------------------------------------------------- | --- | --- | ------------------------------- | ------------------- | ------------------ |
| Showing nodes accounting for 20.30s, 99.64% of 20.37s total |     |     | (*ParallelEngine) tempWorkerRun | runEngine (*Driver) | (*Driver) runAsync |
|                                                             |     |     | 0 of 19.39s (95.17%)            | 0 of 0.79s (3.89%)  | 0 of 0.12s (0.58%) |
Dropped 22 nodes (cum <= 0.10s)
See https://git.io/JfYMW for how to read the graph
|     |     |     |  19.39s                              |  0.79s                                         |  0.12s                           |
| --- | --- | --- | ------------------------------------ | ---------------------------------------------- | -------------------------------- |
|     |     |     | ( * C o m e p m u t u e U            | n i t) ( * P a r a l s l i e m l E n g in      | e) s r e u l n e t c i t m g e o |
|     |     |     | 0 o f   19 H . 3 a 9 n s   d ( l 9 e | 5. 1 7 %) 0   o f  0 . 7 R 9 s u   n (3 . 89 % | ) 0.1 2 s   ( 0 .5 8 %)          |
|     |     |     |  19.39s                              |  0.79s                                         |                                  |
|     |     |     | emu                                  | sim                                            |                                  |
(*ComputeUnit) (*ParallelEngine)
|     |     |     | 0 of 19.39s (95.17%) runEmulation | 0 of 0.79s (3.89%) runRound |     |
| --- | --- | --- | --------------------------------- | --------------------------- | --- |
|     |     |     |  19.39s                           |  0.79s                      |     |
|     |     |     | e m u                             | s y nc                      |     |
|     |     |     | (*Com p ut eUnit)                 | (*Wa i tG roup)             |     |
|     |     |     | 0 of 19.39s (95.17%) runWG        | Wait                        |     |
0.79s (3.89%)
 19.39s
(*ComputeUnit) emu
runWfUntilBarrier
0 of 19.39s (95.17%)
 15.23s
emu
(*ComputeUnit)
|     |     |     | 0 of 15.23s (74.77%) executeInst |     |     |
| --- | --- | --- | -------------------------------- | --- | --- |
 15.23s
emu
(*ALUImpl) Run
0 of 15.23s (74.77%)
|     |  4.16s |                |  0.57s  14.66s |     |     |
| --- | ------ | -------------- | -------------- | --- | --- |
|     |        | (*ALUImpl) emu | (*ALUImpl) emu |     |     |
|     |        | runSMEM        | runFlat        |     |     |
0 of 0.57s (2.82%) 0 of 14.66s (71.95%)
|     |  0.14s                                                              |  0.29s  0.14s                          |  7.45s                               |  7.21s                                |     |
| --- | ------------------------------------------------------------------- | -------------------------------------- | ------------------------------------ | ------------------------------------- | --- |
|     | emu emu                                                             | emu                                    | emu                                  | emu                                   |     |
|     | (*ALUImpl) (*ALUImpl)                                               | (*ALUImpl)                             | (*ALUImpl)                           | (*ALUImpl)                            |     |
|     | runSLOADDWORDX4 0 of 0.14s (0.71%) runSLOADDWORD 0 of 0.29s (1.44%) | runSLOADDWORDX2 0 of 0.14s (0.66%)     | runFlatLoadDWord 0 of 7.45s (36.54%) | runFlatStoreDWord 0 of 7.21s (35.41%) |     |
|     |                                                                     |  0.14s  0.29s  0.14s                   |  7.45s                               |  7.21s                                |     |
|     |                                                                     | emu                                    |                                      | emu                                   |     |
|     |                                                                     | (*storageAccessor) Read                |                                      | (*storageAccessor) Write              |     |
|     |                                                                     | 0 of 12.17s (59.76%)                   |                                      | 0 of 7.21s (35.41%)                   |     |
|     |                                                                     |  9.77s                                 |  2.41s                               |  1.45s  5.76s                         |     |
|     |                                                                     | (*Storage) mem                         | (*pageTableImpl) vm                  | (*Storage) mem                        |     |
|     |                                                                     | 0 of 9.77s (47.95%) Read               | 0 of 3.86s (18.94%) Find             | 0 of 5.76s (28.28%) Write             |     |
|     |                                                                     |  9.77s                                 |  3.72s                               |  0.14s  5.76s                         |     |
|     |                                                                     | vm mem                                 |                                      | vm                                    |     |
|     | (*processTable)                                                     | find createOrGetStorageUnit (*Storage) | (*pageTableImpl) getTable            |                                       |     |
|     | 0 of 3.72s (18.26%)                                                 | 0 of 15.53s (76.23%)                   | 0 of 0.14s (0.68%)                   |                                       |     |
|     |                                                                     |  3.72s                                 |  15.53s  0.14s                       |                                       |     |
|     |                                                                     |  (inline)  (inline)                    |  (inline)                            |                                       |     |
sync
(*Mutex)
Lock
19.39s (95.17%)
FigureA.1: ThesynchronizationblockingprofilefromGotracetoolforReLUbenchmarkparallel
emulation.
108

APPENDIXA. FULLSYNCHRONIZATIONBLOCKINGPROFILES
File: relu
Build ID: 274d7db2213f041fe743ab6c82426062884e144a
Type: delay
| Time: Jul 7, 2024 at 5:39pm (EDT)                           | (*ParallelEngine) sim | (*Driver) driver   |
| ----------------------------------------------------------- | --------------------- | ------------------ |
| Showing nodes accounting for 10.81s, 99.89% of 10.82s total | tempWorkerRun         | runEngine          |
| Dropped 89 nodes (cum <= 0.05s)                             | 0 of 10.25s (94.70%)  | 0 of 0.56s (5.13%) |
See https://git.io/JfYMW for how to read the graph
|     |  10.25s                                     |  0.56s                |
| --- | ------------------------------------------- | --------------------- |
|     | sim                                         | sim                   |
|     | (*TickingComponent) Handle                  | (*ParallelEngine) Run |
|     | 0 of 10.25s (94.70%)                        | 0 of 0.56s (5.13%)    |
|     |  10.17s  0.08s                              |  0.56s                |
|     | (*TickScheduler) sim (*LimitNumMsgPort) sim | (*ParallelEngine) sim |
|     | TickLater Send                              | runRound              |
|     | 0 of 10.17s (93.96%) 0 of 0.08s (0.75%)     | 0 of 0.56s (5.13%)    |
|     |  10.17s  0.08s                              |  0.56s                |
s y nc
|     | s i m s im                                               | (*Wa i tG roup) |
| --- | -------------------------------------------------------- | --------------- |
|     | (*Paral Schedule l e lEngine) (*Direct Send C onnection) | Wait            |
|     | 0 of 10.17s (93.94%) 0 of 0.08s (0.73%)                  | 0.56s (5.13%)   |
|     |  10.17s   0                                              | . 0 8 s         |
 ( i n l i n e)
runtime
sync
|     | chanrecv1 (*Mutex) |     |
| --- | ------------------ | --- |
Lock
|     | 10.17s (94.00%) 0.08s (0.76%) |     |
| --- | ----------------------------- | --- |
FigureA.2: ThesynchronizationblockingprofilefromGotracetoolforReLUbenchmarkparallel
timingsimulation.
109

APPENDIXA. FULLSYNCHRONIZATIONBLOCKINGPROFILES
Type: delay
| S h o w in g   n o d e s  a | c c o u n t in g  f o r  1 .2 9ms, 100% of 1.29ms total |                                       |                    |     |     |     |     |     |
| --------------------------- | ------------------------------------------------------- | ------------------------------------- | ------------------ | --- | --- | --- | --- | --- |
| D ro p p ed   9   n o d e   | s  ( c u m   < =  0 . 0 1 m s)                          | (  t* oeP a 1pr .Wa l slom iem rl skE | nr1 gR0 in 0ue %n) |     |     |     |     |     |
fm   e(
0 2 9 )
See https://git.io/JfYMW for how to read the graph
 1.29ms
s i mods
(* T i c k i .Hn g aC n ml e(1 p o n en t)
0  o f   1 2 9 m 0 0 % )
|     |     |  0.81ms                           |  0.11ms                              |  0.06ms  0.01ms                                                                                          |  0.24ms                                  |  0.02ms                                                          |  0.02ms                                       |                                                                                   |
| --- | --- | --------------------------------- | ------------------------------------ | -------------------------------------------------------------------------------------------------------- | ---------------------------------------- | ---------------------------------------------------------------- | --------------------------------------------- | --------------------------------------------------------------------------------- |
|     |     | w ( *r CT i tae cb  kha           | e c )k w ( r* itC e ai a cc          | r kh o eu n) d ( * M e m d TC r aocs mnk  t r o l l e r) ( * 1TT                                         | t lLc b B ) (                            | * C o m pc iucs u tk( e U n it )                                 | ( * R e o r 2Tdr mieo rb B u f f e r)         | ( * D i r e c tT sC i m co skn n e c t io n )                                     |
|     |     | ics                               | 1T                                   | i                                                                                                        | i sk                                     | 4Tm                                                              | cs k                                          | i                                                                                 |
|     |     | 0 of 0 . 8 1 m (                  | 6 3 .05%) 0 of  0 . 1 m              | s   ( 8 . 39%) 0  o f  0 . 06 m ( 4 . 7 5 % ) 0 of 0. 0                                                  | m (1 .08%) 0 o f                         |  0 . 2   1 8 .6 9 %)                                             | 0  o f  0 . 0 ( 1. 8 0 % )                    | 0  o f   0 .0 2 m  ( 1 . 4 6 % )                                                  |
|     |     |  0.81ms                           |  0.11ms                              |  0.06ms                                                                                                  |  0.01ms                                  |  0.24ms                                                          |  0.02ms                                       |  0.02ms                                                                           |
|     |     | w (u8 *n r CP i tai e cb  eha (le | c )k w (u. r* itC e ai a pc          | r eh o eu n) d ( * M e mr0e6 dCsm r aoos mn  nt rd4 o l l e r) ( 0l*o1                                   | T tl kLb Bp(1 ) (                        | * Cr.uo 2nm pc ipu u  et el1iUn8en it )                          | ( * R eu. o 0nr 2Pdr mieo rb Bl(iu fe8 f e r) | ( * D ffi  or 0re c ta sC imods nM ( n 1ae .nc 4yt io n )                         |
|     |     | 0 of 0r . 1 m ps                  | 6 in3 e.05%) 0 of  r0 1n 1P m        | s   l( i8n .e 39%) 0  o f  0 . p ( . 7 5 % ) 0 of 0.                                                     | mo su  .08%) 0 o                         | f  0 4Pm s ( .6 9 %)                                             | 0  o f  r0 ps e  1n. 0 % )                    | 0  o .w0 2 mr 6 % )                                                               |
|     |     |  0.81ms                           |  0.11ms                              |  0.06ms                                                                                                  |  0.01ms                                  |  0.24ms                                                          |  0.01ms  0.01ms                               |  0.02ms                                                                           |
|     |     | w r i t Sae b  aha                | c )k w r i tC e mcaa                 | r o eu n)S. d d r ao mn                                                                                  | t lLb                                    | cuus ul                                                          | r eo mrb                                      | dr eo rb s imm                                                                    |
|     |     | (r8 *u Cn ct (ge                  | e3 0t iocfk  D0 ( i* r 1e tc         | oh r( y8 3t9a%ge) 0f(  io *nf Ma l0iez. m 0e6 TC r sa  nt rs4 oa. lc7 l 5tei% ro)n) h0a nodfl e0T. ( 0r* | 1aT n ssl  Ba(1t)i.o0n8M%i)ss (          | * S c h e d mR  ne r I m p l)                                    | ( * R e 0bo .o0 rt1 dt om BU ( u 0pf f e r)   | ( * R et. oo0 rp1 Dm os Bw( u 0n. f f e r) ( * L im i t N 2R ue cs vM s g P or t) |
|     |     | 0 of 0 . 1 m s                    | 6 .05%) . 1                          | s   m (                                                                                                  | m 0  o                                   | f  0 . 2 4 ( 1 8 .6 9 %)                                         | 0   o f  s . 9 % ) 0                          | o f  0   8 9 % ) 0   of  0 . 0 m (1 . 4 6% )                                      |
|     |     |  0.81ms                           |                                      |  0.11ms  0.06ms                                                                                          |  0.01ms                                  |  0.24ms                                                          |                                               |  0.02ms                                                                           |
|     |     | sSi m                             | w( rd iit 1Tre a                     | r oo ur ny d dCRm raoes mna                                                                              | t lLb                                    | cuns ult                                                         |                                               | s mfCi m                                                                          |
|     |     | ( * b u f f e re 2Td i cs e  k( n | d e r I m p l) * ei                  | cc kt ) 0f(  io *nf Ma 0lei.zm 0e6 t (dr 4To l 7rla5 en% r)s)                                            | 0 of f0e.t( 0c* 1hT B so  Bt(t1 )o.m08%) | 0E  ( ov* faS  l0uc .ah 2te 4ed mI ee( rr1 In8 ma.6lpI9 ln) %st) |                                               | ( * T ic  N0 k io0 nt2 gi y oR  m (e1 pc. ov4 n e n t)                            |
|     |     | 0   o f   0 .8 m                  | 6 3 . 4 4 % ) 0 of  0 . 1 m          | s   ( 8 . 3 9%) .                                                                                        | m                                        |                                                                  |                                               | 0   of . s 6 % )                                                                  |
|     |     |                                   |                                      |  0.02ms  0.09ms                                                                                          |                                          |  0.24ms                                                          |                                               |  0.02ms                                                                           |
|     |     |                                   | w( *r. r od i cit re ea scss         | rtR  oo u er nya. d )d9%) w( 0r*o. rd iit 9ere es acss r Wt oo ur nyi. d )e0%)                           |                                          | ( * feS  0vc .ah 2le 4Sd cu u nl ed( rP1 I 8gm                   | p l)                                          | ( * T  0Tic .0ikc2 Ss icL m ha  e (tde1 u le r )                                  |
|     |     |                                   | 0 of p 0 0 2 em                      | ( 1 6 0 ofp  0c m   ( 6r t7                                                                              |                                          | 0  o mE s   .m6                                                  | 9 %)                                          | 0  o f km s .r4 6 % )                                                             |
|     |     |                                   |                                      |  0.02ms  0.09ms                                                                                          |                                          |  0.24ms                                                          |  0.01ms                                       |  0.01ms                                                                           |
|     |     |                                   | w( r i sit e ac                      | r o u nyM. d w( rd i rit 9ire acMs rt oo ur nys. d                                                       |                                          |                                                                  | c u                                           |                                                                                   |
|     |     |  0.82ms                           | 0c* ed rs eR                         | et ao dr ) *w. ee )  0.06ms                                                                              |  0.01ms                                  | ( * fGS  0Cc h 2oe 4md                                           | pu lle  e (tr 1iIo8 mn.6Mp le)                |                                                                                   |
|     |     |                                   | 0 porfo  . 0 2 m                     | s   ( 1 6 9is%s) 0 of  0 0 tm   (i 6s 7 0%)                                                              |                                          | send0W  o .                                                      | m s 9 %ss)age                                 |                                                                                   |
|     |     |                                   |  0.02ms                              |  0.09ms                                                                                                  |                                          |                                                                  |                                               |                                                                                   |
|     |     |                                   | w(c0 rdF0 iit 2rre moea rt o ur ny d | w( rd iit t e ac r oot( u ny d                                                                           |                                          |                                                                  |                                               |                                                                                   |
|     |     |                                   | 0 foeft  *h. mc Bo 1o .t )to9m%)     | 0 of w * .r 0i re eB ot tr 6o .m )                                                                       |                                          |  0.24ms                                                          |                                               |                                                                                   |
|     |     |                                   | s   ( 6                              | 0 9 m s   7 0%)                                                                                          |                                          |                                                                  |                                               |                                                                                   |
|     |     |                                   |  0.02ms                              |  0.09ms                                                                                                  |                                          |                                                                  |                                               |  0.02ms                                                                           |
( * L i m i tN sue imm M sg P o rt )
0  o f   1. 2 7Sm sn d( 98 .5 4 % )
 1.25ms  0.02ms
|     |     | ( * D i r e c tSm sCe im on  nd( | n e c ti o n) ( *s Ew 02Sni tde cPns | h  doi ni ng t ) |     |     |     |     |
| --- | --- | -------------------------------- | ------------------------------------ | ---------------- | --- | --- | --- | --- |
|     |     | 0   o f  1 . 2 5 s               | 9 7 .0 8 % ) 0 of  0 . m             | ( 1 . 4 6%)      |     |     |     |     |
 1.25ms  0.02ms
sync
(*Mutex)
Lock
1.29ms (100%)
FigureA.3: Thesynchronizationblockingprofileafterapplyingthereversing-orderstrategy.
110
